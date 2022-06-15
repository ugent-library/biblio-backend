package models

import (
	"encoding/json"
	"time"
)

/*
	ID: identifier to use as foreign key. Possibly not unique?
	Type: type of model (user, department)
		Combination of id,type must be unique
*/
type RelatedObject struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Data        json.RawMessage `json:"data"`
	DateCreated *time.Time      `json:"date_created,omitempty"`
	DateUpdated *time.Time      `json:"date_updated,omitempty"`
}

func (ro *RelatedObject) setData(data interface{}) error {
	b, e := json.Marshal(data)
	if e != nil {
		return e
	}
	ro.Data = b
	return nil
}
