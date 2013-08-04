package application

import (
	"github.com/emicklei/landskape/model"
	"testing"
)

func TestDotBuilderOneConnection(t *testing.T) {
	c := []model.Connection{model.Connection{From: "A", To: "B", Type: "T"}}
	b := NewDotBuilder()
	b.BuildFromAll(c)
	b.WriteDotFile("/tmp/TestDotBuilderOneConnection.dot")
}
