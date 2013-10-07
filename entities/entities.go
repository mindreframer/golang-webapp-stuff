// Package entities describes the domain entities for the openvoicedata.org project
package entities

import (
	"fmt"
)

type ErrNotFound struct {
	Entity interface{}
	Id     int
}

func NewErrNotFound(entity interface{}, id int) ErrNotFound {
	return ErrNotFound{entity, id}
}

func (nf ErrNotFound) Error() string {
	return fmt.Sprintf("%T id:%d not found", nf.Entity, nf.Id)
}
