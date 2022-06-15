package related_object_store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/ulid"
)

/*
	TODO: add context
		  add transactions (even needed as this has its own db connection?)
*/

type Store struct {
	db     *pgxpool.Pool
	source backends.RelatedObjectSource
}

/*
	ros := related_object_store.New(
		"postgres://biblio:biblio@localhost:5432/biblio_backend?sslmode=disable",
		biblioClient,
	)
*/
func New(dsn string, source backends.RelatedObjectSource) (*Store, error) {

	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to postgres")
	}

	return &Store{
		source: source,
		db:     db,
	}, nil

}

/*
	fetch object from store
	id possibly not unique, so we require "type" too
	TODO: cache
*/
func (ros *Store) Get(t string, id string) (*models.RelatedObject, error) {
	ro := models.RelatedObject{}

	sql := `select id,date_created,date_updated,type,data
	from related_objects
	where id = $1 and type = $2 limit 1`
	err := ros.db.QueryRow(
		context.Background(),
		sql,
		id,
		t).
		Scan(&ro.ID, &ro.DateCreated, &ro.DateUpdated, &ro.Type, &ro.Data)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to fetch row")
	}

	return &ro, nil
}

/*
	TODO: a query like "where in" doesn't necessarily
		  fetch all records in that list
		  Is it an error if there are less records
		  than ids given? Or is that more a problem
		  for the client code?
*/
func (ros *Store) GetMultiple(t string, ids []string) ([]*models.RelatedObject, error) {
	objects := make([]*models.RelatedObject, 0, 100)

	pgIds := &pgtype.TextArray{}
	pgIds.Set(ids)
	sql := `select id,date_created,date_updated,type,data
	from related_objects
	where id = any($1) and type = $2`

	c, err := ros.db.Query(context.Background(), sql, pgIds, t)
	if err != nil {
		return objects, err
	}
	defer c.Close()

	for c.Next() {
		ro := models.RelatedObject{}
		e := c.Scan(&ro.ID, &ro.DateCreated, &ro.DateUpdated, &ro.Type, &ro.Data)
		if e != nil {
			return objects, e
		}
		objects = append(objects, &ro)
	}

	return objects, nil
}

/*
	fetch data from url
	and update object
*/
func (ros *Store) Sync(ro *models.RelatedObject) error {

	switch ro.Type {
	case "user":
		u, e := ros.source.GetPerson(ro.ID) //GetUser only returns active users
		if e != nil {
			return e
		}
		bytes, e := json.Marshal(u)
		if e != nil {
			return e
		}
		ro.Data = bytes
	case "department":
		d, e := ros.source.GetOrganization(ro.ID)
		if e != nil {
			return e
		}
		bytes, e := json.Marshal(d)
		if e != nil {
			return e
		}
		ro.Data = bytes
	case "project":
		p, e := ros.source.GetProject(ro.ID)
		if e != nil {
			return e
		}
		bytes, e := json.Marshal(p)
		if e != nil {
			return e
		}
		ro.Data = bytes
	default:
		return errors.New(fmt.Sprintf("Unable to mapping type %s", ro.Type))
	}

	return nil
}

/*
	save back to table
*/
func (ros *Store) Save(ro *models.RelatedObject) error {
	now := time.Now()
	ro.DateUpdated = &now
	if ro.DateCreated == nil {
		ro.DateCreated = &now
	}
	if ro.ID == "" {
		ro.ID, _ = ulid.Generate()
	}
	sql := `insert into
	related_objects(id,date_created,date_updated,type,data) values($1,$2,$3,$4,$5)
	on conflict(id,type)
	do update
	set
		date_created = excluded.date_created,
		date_updated = excluded.date_updated,
		data = excluded.data
	`
	_, err := ros.db.Exec(
		context.Background(),
		sql,
		ro.ID,
		ro.DateCreated,
		ro.DateUpdated,
		ro.Type,
		ro.Data)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Unable to save related_object %s", ro.ID))
	}

	return nil
}

func (ros *Store) SyncAndSave(ro *models.RelatedObject) error {
	if e := ros.Sync(ro); e != nil {
		return e
	}
	if e := ros.Save(ro); e != nil {
		return e
	}
	return nil
}

func (ros *Store) Each(fn func(ro *models.RelatedObject)) error {
	sql := "select id,date_created,date_updated,type,data from related_objects"
	c, err := ros.db.Query(context.Background(), sql)
	if err != nil {
		return err
	}
	defer c.Close()

	for c.Next() {
		ro := models.RelatedObject{}
		e := c.Scan(&ro.ID, &ro.DateCreated, &ro.DateUpdated, &ro.Type, &ro.Data)
		if e != nil {
			return e
		}
		fn(&ro)
	}

	return nil
}

func (ros *Store) SyncAndSaveAll() error {
	objects := make([]*models.RelatedObject, 0, 100)

	// collect data (do not update while iterating)
	errorEach := ros.Each(func(ro *models.RelatedObject) {
		objects = append(objects, ro)
	})

	if errorEach != nil {
		return errorEach
	}

	// synchronize and update
	for _, ro := range objects {
		e := ros.SyncAndSave(ro)
		if e != nil {
			return e
		}
	}

	return nil
}
