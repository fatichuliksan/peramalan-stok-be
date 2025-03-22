package general

import (
	"peramalan-stok-be/src/helper/postgre"
)

type generalHelper struct {
	DB postgre.Database
}

// Interface ...
type Interface interface {
	ContainString(a string, list []string) bool
}

// NewGeneralHelper ...
func NewGeneralHelper(db postgre.Database) Interface {
	return &generalHelper{
		DB: db,
	}
}

// ContainString ...
func (t *generalHelper) ContainString(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
