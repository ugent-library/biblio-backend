package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/context"
	"github.com/unrolled/render"
)

type DatasetPublications struct {
	Base
	store                    backends.Repository
	publicationSearchService backends.PublicationSearchService
}

func NewDatasetPublications(base Base, store backends.Repository,
	publicationSearchService backends.PublicationSearchService) *DatasetPublications {
	return &DatasetPublications{
		Base:                     base,
		store:                    store,
		publicationSearchService: publicationSearchService,
	}
}

func (c *DatasetPublications) Choose(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	datasetPubs, err := c.store.GetDatasetPublications(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	datasetPubIDs := make([]string, len(datasetPubs))
	for i, d := range datasetPubs {
		datasetPubIDs[i] = d.ID
	}

	searchArgs := models.NewSearchArgs()
	// empty es exclude filter leads to empty results
	if len(datasetPubIDs) > 0 {
		searchArgs.Filters["!id"] = datasetPubIDs
	}

	hits, err := c.userPublications(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/publications/_modal", c.ViewData(r, struct {
		Dataset *models.Dataset
		Hits    *models.PublicationHits
	}{
		dataset,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetPublications) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datasetPubs, err := c.store.GetDatasetPublications(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	datasetPubIDs := make([]string, len(datasetPubs))
	for i, d := range datasetPubs {
		datasetPubIDs[i] = d.ID
	}

	searchArgs := models.NewSearchArgs()
	searchArgs.Query = r.Form["search"][0]
	// empty es exclude filter leads to empty results
	if len(datasetPubIDs) > 0 {
		searchArgs.Filters["exclude"] = datasetPubIDs
	}

	hits, err := c.userPublications(context.GetUser(r.Context()).ID, searchArgs)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/publications/_modal_hits", c.ViewData(r, struct {
		Dataset *models.Dataset
		Hits    *models.PublicationHits
	}{
		dataset,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetPublications) Add(w http.ResponseWriter, r *http.Request) {
	pubID := mux.Vars(r)["publication_id"]

	dataset := context.GetDataset(r.Context())

	pub, err := c.store.GetPublication(pubID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = c.store.AddPublicationDataset(pub, dataset); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	publications, err := c.store.GetDatasetPublications(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/publications/_show", c.ViewData(r, struct {
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
	}{
		dataset,
		publications,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetPublications) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	pubID := mux.Vars(r)["publication_id"]

	c.Render.HTML(w, http.StatusOK, "dataset/publications/_modal_confirm_removal", c.ViewData(r, struct {
		DatasetID     string
		PublicationID string
	}{
		id,
		pubID,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetPublications) Remove(w http.ResponseWriter, r *http.Request) {
	pubID := mux.Vars(r)["publication_id"]

	dataset := context.GetDataset(r.Context())

	pub, err := c.store.GetPublication(pubID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = c.store.RemovePublicationDataset(pub, dataset); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	publications, err := c.store.GetDatasetPublications(dataset)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/publications/_show", c.ViewData(r, struct {
		Dataset             *models.Dataset
		DatasetPublications []*models.Publication
	}{
		dataset,
		publications,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetPublications) userPublications(userID string, args *models.SearchArgs) (*models.PublicationHits, error) {
	searcher := c.publicationSearchService.WithScope("status", "private", "public")

	switch args.FilterFor("scope") {
	case "created":
		searcher = searcher.WithScope("creator_id", userID)
	case "contributed":
		searcher = searcher.WithScope("author.id", userID)
	default:
		searcher = searcher.WithScope("creator_id|author.id", userID)
	}
	delete(args.Filters, "scope")

	return searcher.Search(args)
}