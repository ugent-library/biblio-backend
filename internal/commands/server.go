package commands

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/services"
	"github.com/ugent-library/biblio-backend/services/webapp"
)

func init() {
	serverCmd.PersistentFlags().String("base-url", "", "base url")

	serverStartCmd.Flags().String("mode", defaultMode, "server mode (development, staging or production)")
	serverStartCmd.Flags().String("host", defaultHost, "server host")
	serverStartCmd.Flags().Int("port", defaultPort, "server port")
	serverStartCmd.Flags().String("session-name", defaultSessionName, "session name")
	serverStartCmd.Flags().String("session-secret", "", "session secret")
	serverStartCmd.Flags().Int("session-max-age", defaultSessionMaxAge, "session lifetime")
	serverStartCmd.Flags().String("csrf-name", "", "csrf cookie name")
	serverStartCmd.Flags().String("csrf-secret", "", "csrf cookie secret")

	webapp.AddCommands(serverCmd, Services())
	serverCmd.AddCommand(serverStartCmd)
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server [command]",
	Short: "The biblio-backend HTTP server",
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "start the http server",
	Run: func(cmd *cobra.Command, args []string) {
		e := Services()
		e.Store.AddPublicationListener(func(p *models.Publication) {
			cPub, cE := preprocessPublicationForIndex(p)
			if cE != nil {
				log.Println(fmt.Errorf("error indexing publication %s: %w", p.ID, cE))
				return
			}
			if err := e.PublicationSearchService.Index(cPub); err != nil {
				log.Println(fmt.Errorf("error indexing publication %s: %w", p.ID, err))
			}
		})
		e.Store.AddDatasetListener(func(d *models.Dataset) {
			if err := e.DatasetSearchService.Index(d); err != nil {
				log.Println(fmt.Errorf("error indexing dataset %s: %w", d.ID, err))
			}
		})

		wa, err := webapp.New(e)
		if err != nil {
			log.Fatal(err)
		}

		if err = services.Start(wa); err != nil {
			log.Fatal(err)
		}
	},
}

func preprocessPublicationForIndex(p *models.Publication) (*models.Publication, error) {
	pub := p.Clone()

	log.Println("preprocessPublicationForIndex")

	//update contributors
	if e := updateContributors(pub); e != nil {
		return nil, errors.Wrapf(
			e,
			"Unable to update contributors from related objects")
	}

	//update departments
	if e := updateDepartments(pub); e != nil {
		return nil, errors.Wrapf(
			e,
			"Unable to update departments from related objects")
	}

	//update projects
	if e := updateProjects(pub); e != nil {
		return nil, errors.Wrapf(
			e,
			"Unable to update projects from related objects")
	}

	return pub, nil
}

func updateProjects(pub *models.Publication) error {
	log.Println("updateProjects")
	services := Services()
	ids := make([]string, 0, 10)
	for _, pr := range pub.Project {
		if pr.ID != "" {
			ids = append(ids, pr.ID)
		}
	}
	ids = validation.ArrayUniq(ids)
	if len(ids) == 0 {
		return nil
	}
	records, e := services.RelatedObjectService.GetMultiple("project", ids)
	if e != nil {
		return e
	}
	if len(records) == 0 {
		return nil
	}
	oldProjects := make([]*models.PublicationProject, len(records))
	for i, record := range records {
		oldPr := models.PublicationProject{}
		if e := json.Unmarshal(record.Data, &oldPr); e != nil {
			return errors.Wrap(e,
				fmt.Sprintf(
					"unable to decode project from related_object %s",
					record.ID,
				),
			)
		}
		oldPr.ID = record.ID
		oldProjects[i] = &oldPr
	}
	for _, project := range pub.Project {
		if project.ID == "" {
			continue
		}
		for _, oldProject := range oldProjects {
			if oldProject.ID == project.ID {
				project.Name = oldProject.Name
			}
		}
	}
	return nil
}

func updateDepartments(pub *models.Publication) error {
	log.Println("updateDepartments")
	services := Services()
	ids := make([]string, 0, 10)
	for _, dep := range pub.Department {
		if dep.ID != "" {
			ids = append(ids, dep.ID)
		}
	}
	ids = validation.ArrayUniq(ids)
	if len(ids) == 0 {
		return nil
	}
	records, e := services.RelatedObjectService.GetMultiple("department", ids)
	if e != nil {
		return e
	}
	if len(records) == 0 {
		return nil
	}
	oldDeps := make([]*models.PublicationDepartment, len(records))
	for i, record := range records {
		oldDep := models.PublicationDepartment{}
		if e := json.Unmarshal(record.Data, &oldDep); e != nil {
			return errors.Wrap(e,
				fmt.Sprintf(
					"unable to decode department from related_object %s",
					record.ID,
				),
			)
		}
		oldDep.ID = record.ID
		oldDeps[i] = &oldDep
	}
	for _, dep := range pub.Department {
		for _, oldDep := range oldDeps {
			if dep.ID != oldDep.ID {
				continue
			}
			dep.Tree = oldDep.Clone().Tree
		}
	}
	return nil
}

func updateContributors(pub *models.Publication) error {
	log.Println("updateContributors")
	services := Services()
	// update contributors
	{
		ids := extractUserIds(pub.Author...)
		ids = append(ids, extractUserIds(pub.Editor...)...)
		ids = append(ids, extractUserIds(pub.Supervisor...)...)
		ids = validation.ArrayUniq(ids)
		log.Printf("found %d userIds in pub %s", len(ids), pub.ID)
		if len(ids) == 0 {
			return nil
		}

		records, e := services.RelatedObjectService.GetMultiple("user", ids)
		if e != nil {
			return e
		}
		log.Printf("found %d users in related for pub %s", len(records), pub.ID)
		if len(records) == 0 {
			return nil
		}

		oldUsers := make([]*models.Contributor, len(records))

		for i, record := range records {
			c := models.Contributor{}
			if e := json.Unmarshal(record.Data, &c); e != nil {
				return errors.Wrap(e,
					fmt.Sprintf(
						"unable to decode user from related_object %s",
						record.ID,
					),
				)
			}
			c.ID = record.ID
			oldUsers[i] = &c
		}

		for _, c := range pub.Author {
			if c.ID == "" {
				continue
			}
			for _, oldUser := range oldUsers {
				if oldUser.ID == c.ID {
					copyContributor(c, oldUser)
				}
			}
		}
		for _, c := range pub.Editor {
			if c.ID == "" {
				continue
			}
			for _, oldUser := range oldUsers {
				if oldUser.ID == c.ID {
					copyContributor(c, oldUser)
				}
			}
		}
		for _, c := range pub.Supervisor {
			if c.ID == "" {
				continue
			}
			for _, oldUser := range oldUsers {
				if oldUser.ID == c.ID {
					copyContributor(c, oldUser)
				}
			}
		}
	}
	return nil
}

func extractUserIds(contributors ...*models.Contributor) []string {
	userIds := make([]string, 0, 10)
	for _, contributor := range contributors {
		if contributor.ID == "" {
			continue
		}
		userIds = append(userIds, contributor.ID)
	}
	return userIds
}

func copyContributor(dest *models.Contributor, src *models.Contributor) {
	dest.FirstName = src.FirstName
	dest.LastName = src.LastName
	dest.FullName = src.FullName
	dest.CreditRole = nil
	dest.CreditRole = append(dest.CreditRole, src.CreditRole...)
	dest.ORCID = src.ORCID
	dest.UGentID = nil
	dest.UGentID = append(dest.UGentID, src.UGentID...)
}
