package webservice

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/landskape/application"
	"github.com/emicklei/landskape/model"
	"net/http"
)

const (
	NO_UPDATE = false
)

type SystemResource struct {
	Logic application.Logic
}

func (s SystemResource) Register() {
	ws := new(restful.WebService)
	ws.Path("/{scope}/systems").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_XML, restful.MIME_JSON).
		Param(ws.PathParameter("scope", "organization name to group system and connections").DataType("string"))

	idParam := ws.PathParameter("id", "identifier of the system").DataType("string")

	ws.Route(ws.GET("").To(s.getAll).
		// docs
		Doc("list all known systems").
		Writes([]model.System{}))

	ws.Route(ws.GET("/{id}").To(s.get).
		// docs
		Doc("get the system using its id").
		Param(idParam).
		Writes(model.System{})) // to the response

	ws.Route(ws.PUT("/{id}").To(s.put).
		// docs
		Doc("create the system using its id").
		Param(idParam).
		Reads(model.System{})) // from the request

	ws.Route(ws.POST("").To(s.post).
		// docs
		Doc("update the system using its id").
		Param(idParam).
		Reads(model.System{})) // from the request

	ws.Route(ws.DELETE("/{id}").To(s.delete).
		// docs
		Doc("delete the system using its id").
		Param(idParam))

	restful.Add(ws)
}

// DELETE /{scope}/systems/{id}
func (s SystemResource) delete(req *restful.Request, resp *restful.Response) {
	scope := req.PathParameter("scope")
	id := req.PathParameter("id")
	err := s.Logic.DeleteSystem(scope, id)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
}

// GET /{scope}/systems/{id}
func (s SystemResource) get(req *restful.Request, resp *restful.Response) {
	scope := req.PathParameter("scope")
	id := req.PathParameter("id")
	app, err := s.Logic.GetSystem(scope, id)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(app)
}

func (s SystemResource) getAll(req *restful.Request, resp *restful.Response) {
	scope := req.PathParameter("scope")
	apps, err := s.Logic.AllSystems(scope)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteEntity(apps)
}

// POST /{scope}/systems/
func (s SystemResource) post(req *restful.Request, resp *restful.Response) {
	app := new(model.System)
	err := req.ReadEntity(app)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if app.Id == "" {
		err := restful.NewError(model.MISMATCH_ID, "Id is missing")
		resp.WriteServiceError(http.StatusBadRequest, err)
		return
	}
	_, err = s.Logic.SaveSystem(app)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
	}
}

// PUT /{scope}/systems/{id}
func (s SystemResource) put(req *restful.Request, resp *restful.Response) {
	scope := req.PathParameter("scope")
	id := req.PathParameter("id")
	app := new(model.System)
	err := req.ReadEntity(app)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	if app.Id != id {
		err := restful.NewError(model.MISMATCH_ID, fmt.Sprintf("Id mismatch: %v != %v", app.Id, id))
		resp.WriteServiceError(http.StatusBadRequest, err)
		return
	}
	if app.Scope != "" && app.Scope != scope {
		err := restful.NewError(model.MISMATCH_SCOPE, fmt.Sprintf("Scope mismatch: %v != %v", app.Scope, scope))
		resp.WriteServiceError(http.StatusBadRequest, err)
		return
	}
	if s.Logic.ExistsSystem(scope, id) {
		err := restful.NewError(http.StatusConflict, "System already exists:"+id)
		resp.WriteServiceError(http.StatusConflict, err)
		return
	}
	_, err = s.Logic.SaveSystem(app)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	resp.WriteHeader(http.StatusCreated)
}
