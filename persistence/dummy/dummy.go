// Package dummy implements dummy storage for API entities
package dummy

import (
	"github.com/ku-ovdp/api/persistence"
)

type dummyBackend struct{}

func (d dummyBackend) Init(args ...interface{}) {}

func init() {
	persistence.Register("dummy", dummyBackend{})
}
