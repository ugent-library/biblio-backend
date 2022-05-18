package engine

import (
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/backends/filestore"
	"github.com/ugent-library/biblio-backend/internal/tasks"
	"github.com/ugent-library/go-orcid/orcid"
)

type Config struct {
	ORCIDSandbox             bool
	ORCIDClient              *orcid.MemberClient
	Store                    backends.Store
	FileStore                *filestore.Store
	DatasetSearchService     backends.DatasetSearchService
	PublicationSearchService backends.PublicationSearchService
	backends.PersonService
	backends.ProjectService
	backends.UserService
	backends.OrganizationSearchService
	backends.PersonSearchService
	backends.ProjectSearchService
	backends.LicenseSearchService
	backends.MediaTypeSearchService
	PublicationSources  map[string]backends.PublicationGetter
	DatasetSources      map[string]backends.DatasetGetter
	PublicationEncoders map[string]backends.PublicationEncoder
	PublicationDecoders map[string]backends.PublicationDecoderFactory
	Tasks               *tasks.Hub
}

type Engine struct {
	Config
}

func New(c Config) (*Engine, error) {
	e := &Engine{c}

	return e, nil
}
