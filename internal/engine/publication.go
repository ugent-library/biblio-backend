package engine

import (
	"fmt"

	"github.com/go-playground/form/v4"
)

type Publication map[string]interface{}

type PublicationHits struct {
	Total        int           `json:"total"`
	Page         int           `json:"page"`
	LastPage     int           `json:"last_page"`
	PreviousPage bool          `json:"previous_page"`
	NextPage     bool          `json:"next_page"`
	Hits         []Publication `json:"hits"`
}

func (e *Engine) UserPublications(userID string, q *Query) (*PublicationHits, error) {
	qp, err := form.NewEncoder().Encode(q)
	if err != nil {
		return nil, err
	}
	hits := &PublicationHits{}
	if _, err := e.get(fmt.Sprintf("/user/%s/publication", userID), qp, hits); err != nil {
		return nil, err
	}
	return hits, nil
}

func (e *Engine) GetPublication(id string) (*Publication, error) {
	pub := &Publication{}
	if _, err := e.get(fmt.Sprintf("/publication/%s", id), nil, pub); err != nil {
		return nil, err
	}
	return pub, nil
}
