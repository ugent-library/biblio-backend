package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ugent-library/biblio-backend/internal/backends/biblio"
	"github.com/ugent-library/biblio-backend/internal/models"
	"github.com/ugent-library/biblio-backend/internal/related_object_store"
	"github.com/ugent-library/biblio-backend/internal/validation"
)

func init() {
	relatedObjectCmd.AddCommand(relatedObjectExtractCmd)
	rootCmd.AddCommand(relatedObjectCmd)
}

func newBiblioClient() *biblio.Client {
	return biblio.New(biblio.Config{
		URL:      viper.GetString("frontend-url"),
		Username: viper.GetString("frontend-username"),
		Password: viper.GetString("frontend-password"),
	})
}

func newRelatedObjectStore() *related_object_store.Store {
	store, e := related_object_store.New(
		viper.GetString("pg-conn"),
		newBiblioClient(),
	)
	if e != nil {
		log.Fatalln(e)
	}
	return store
}

var relatedObjectCmd = &cobra.Command{
	Use:   "related_object [command]",
	Short: "Related object commands",
}

var relatedObjectExtractCmd = &cobra.Command{
	Use:   "extract",
	Short: "extract related objects from existing users, projects and departments",
	Run: func(cmd *cobra.Command, args []string) {

		store := newStore()
		ros := newRelatedObjectStore()
		userIdsNotFound := make([]string, 0, 100)
		depIdsNotFound := make([]string, 0, 100)
		projectIdsNotFound := make([]string, 0, 100)

		/*
			TODO:
			- also related objects hidden in older (unused) snapshots?
			- only execute this during migration
			- also do this for (new) datasets created in LibreCat
			  and migrated to here
		*/
		store.EachPublicationSnapshot(func(pub *models.Publication) bool {

			// extract and store projects
			projectIds := make([]string, 0, 10)
			for _, p := range pub.Project {
				if p.ID == "" {
					continue
				}
				projectIds = append(projectIds, p.ID)
			}
			oldProjects := make([]*models.RelatedObject, 0, 10)
			if len(projectIds) > 0 {
				var e error
				oldProjects, e = ros.GetMultiple("project", projectIds)
				if e != nil {
					log.Fatalln(e)
				}
			}
			for _, pId := range projectIds {
				// already inserted
				var found bool = false
				for _, p := range oldProjects {
					if p.ID == pId {
						found = true
						break
					}
				}
				if found {
					continue
				}
				// already tried to fetch from biblio
				if validation.InArray(projectIdsNotFound, pId) {
					continue
				}
				ro := models.RelatedObject{
					Type: "project",
					ID:   pId,
				}
				e := ros.SyncAndSave(&ro)
				if e != nil {
					if _, ok := e.(models.HttpNotFound); ok {
						projectIdsNotFound = append(projectIdsNotFound, pId)
						log.Printf(
							"publication %s [snapshot: %s] has non existing project %s",
							pub.ID,
							pub.SnapshotID,
							pId,
						)
						continue
					}
					log.Printf("Unable to fetch/save project %s : %s", pId, e.Error())
				} else {
					log.Printf(
						"Add project[%s] for %s [snapshot: %s]",
						pId,
						pub.ID,
						pub.SnapshotID,
					)
				}
			}

			// extract and store users
			userIds := make([]string, 0, 10)
			for _, c := range pub.Author {
				if c.ID == "" {
					continue
				}
				if validation.InArray(userIds, c.ID) {
					continue
				}
				userIds = append(userIds, c.ID)
			}
			for _, c := range pub.Editor {
				if c.ID == "" {
					continue
				}
				if validation.InArray(userIds, c.ID) {
					continue
				}
				userIds = append(userIds, c.ID)
			}
			for _, c := range pub.Supervisor {
				if c.ID == "" {
					continue
				}
				if validation.InArray(userIds, c.ID) {
					continue
				}
				userIds = append(userIds, c.ID)
			}
			oldUsers := make([]*models.RelatedObject, 0, 10)
			if len(userIds) > 0 {
				var e error
				oldUsers, e = ros.GetMultiple("user", userIds)
				if e != nil {
					log.Fatalln(e)
				}
			}
			for _, userId := range userIds {
				var found bool = false
				for _, u := range oldUsers {
					if u.ID == userId {
						found = true
						break
					}
				}
				if found {
					continue
				}
				if validation.InArray(userIdsNotFound, userId) {
					continue
				}
				ro := models.RelatedObject{
					Type: "user",
					ID:   userId,
				}
				e := ros.SyncAndSave(&ro)
				if e != nil {
					if _, ok := e.(models.HttpNotFound); ok {
						userIdsNotFound = append(userIdsNotFound, userId)
						log.Printf(
							"publication %s [snapshot: %s] has non existing user %s",
							pub.ID,
							pub.SnapshotID,
							userId,
						)
						continue
					}
					log.Printf("Unable to fetch/save user %s : %s", userId, e.Error())
				} else {
					log.Printf(
						"Add user[%s], found in publication %s [snapshot: %s]",
						userId,
						pub.ID,
						pub.SnapshotID,
					)
				}
			}

			// extract and store departments
			depIds := make([]string, 0, 10)
			for _, dep := range pub.Department {
				if dep.ID == "" {
					continue
				}
				if validation.InArray(depIds, dep.ID) {
					continue
				}
				depIds = append(depIds, dep.ID)
			}
			oldDepartments := make([]*models.RelatedObject, 0, 10)
			if len(depIds) > 0 {
				var e error
				oldDepartments, e = ros.GetMultiple("department", depIds)
				if e != nil {
					log.Fatalln(e)
				}
			}
			for _, depId := range depIds {
				var found bool = false
				for _, dep := range oldDepartments {
					if dep.ID == depId {
						found = true
						break
					}
				}
				if found {
					continue
				}
				if validation.InArray(depIdsNotFound, depId) {
					continue
				}
				ro := models.RelatedObject{
					Type: "department",
					ID:   depId,
				}
				e := ros.SyncAndSave(&ro)
				if e != nil {
					if _, ok := e.(models.HttpNotFound); ok {
						depIdsNotFound = append(depIdsNotFound, depId)
						log.Printf(
							"publication %s [snapshot: %s] has non existing user %s",
							pub.ID,
							pub.SnapshotID,
							depId,
						)
						continue
					}
					log.Printf("Unable to fetch/save department %s : %s", depId, e.Error())
				} else {
					log.Printf(
						"Add department[%s] for %s [snapshot: %s]",
						depId,
						pub.ID,
						pub.SnapshotID,
					)
				}
			}
			return true
		})

		// REPORT
		log.Printf("projects not found: %d", len(projectIdsNotFound))
		log.Printf("users not found: %d", len(userIdsNotFound))
		log.Printf("departments not found: %d", len(depIdsNotFound))

	},
}
