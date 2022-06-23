package controllers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/ugent-library/biblio-backend/internal/backends"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/services/webapp/context"
	"github.com/unrolled/render"
)

type DatasetDepartments struct {
	Base
	store                     backends.Repository
	organizationSearchService backends.OrganizationSearchService
	organizationService       backends.OrganizationService
}

func NewDatasetDepartments(base Base, store backends.Repository,
	organizationSearchService backends.OrganizationSearchService,
	organizationService backends.OrganizationService) *DatasetDepartments {
	return &DatasetDepartments{
		Base:                      base,
		store:                     store,
		organizationSearchService: organizationSearchService,
		organizationService:       organizationService,
	}
}

func (c *DatasetDepartments) List(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	// Get 20 random departments (no search, init state)
	hits, _ := c.organizationSearchService.SuggestOrganizations("")

	c.Render.HTML(w, http.StatusOK, "dataset/departments/_modal", c.ViewData(r, struct {
		Dataset *models.Dataset
		Hits    []models.Completion
	}{
		dataset,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDepartments) ActiveSearch(w http.ResponseWriter, r *http.Request) {
	dataset := context.GetDataset(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get 20 results from the search query
	query := r.Form["search"]
	hits, _ := c.organizationSearchService.SuggestOrganizations(query[0])

	c.Render.HTML(w, http.StatusOK, "dataset/departments/_modal_hits", c.ViewData(r, struct {
		Dataset *models.Dataset
		Hits    []models.Completion
	}{
		dataset,
		hits,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDepartments) Add(w http.ResponseWriter, r *http.Request) {
	departmentID := mux.Vars(r)["department_id"]

	// because mux var {department_id} is not properly decoded
	// i.e. LW06%2A remains LW06%2A instead of LW06*
	departmentID, _ = url.PathUnescape(departmentID)
	dataset := context.GetDataset(r.Context())

	org, err := c.organizationService.GetOrganization(departmentID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	datasetDepartment := models.DatasetDepartment{
		ID: departmentID,
	}
	for _, o := range org.Tree {
		datasetDepartment.Tree = append(datasetDepartment.Tree, models.DatasetDepartmentRef{ID: o.ID})
	}
	dataset.Department = append(dataset.Department, datasetDepartment)
	savedDataset := dataset.Clone()

	if err := c.store.SaveDataset(dataset); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/departments/_show", c.ViewData(r, struct {
		Dataset *models.Dataset
	}{
		savedDataset,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDepartments) ConfirmRemove(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	departmentId := mux.Vars(r)["department_id"]
	// because mux var {department_id} is not properly decoded
	// i.e. LW06%2A remains LW06%2A instead of LW06*
	departmentId, _ = url.PathUnescape(departmentId)

	c.Render.HTML(w, http.StatusOK, "dataset/departments/_modal_confirm_removal", c.ViewData(r, struct {
		ID           string
		DepartmentID string
	}{
		id,
		departmentId,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}

func (c *DatasetDepartments) Remove(w http.ResponseWriter, r *http.Request) {
	departmentId := mux.Vars(r)["department_id"]
	// because mux var {department_id} is not properly decoded
	// i.e. LW06%2A remains LW06%2A instead of LW06*
	departmentId, _ = url.PathUnescape(departmentId)

	dataset := context.GetDataset(r.Context())

	departments := make([]models.DatasetDepartment, len(dataset.Department))
	copy(departments, dataset.Department)

	var removeKey int = -1
	for key, department := range departments {
		if department.ID == departmentId {
			removeKey = key
		}
	}
	if removeKey >= 0 {
		departments = append(departments[:removeKey], departments[removeKey+1:]...)
	}
	dataset.Department = departments

	savedDataset := dataset.Clone()
	err := c.store.SaveDataset(dataset)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.Render.HTML(w, http.StatusOK, "dataset/departments/_show", c.ViewData(r, struct {
		Dataset *models.Dataset
	}{
		savedDataset,
	}),
		render.HTMLOptions{Layout: "layouts/htmx"},
	)
}