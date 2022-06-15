package models

type HttpNotFound struct {
	Message string
}

func (e HttpNotFound) Error() string {
	return e.Message
}
