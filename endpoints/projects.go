// Package endpoints describes the api endpoints
package endpoints

import (
	. "github.com/ku-ovdp/api/entities"
	. "github.com/ku-ovdp/api/repository"
	"github.com/ku-ovdp/api/stats"
	"github.com/traviscline/go-restful"
	"net/http"
	"strconv"
)

type projectService struct {
	*restful.WebService
	repository ProjectRepository
}

func NewProjectService(apiRoot string, repository ProjectRepository) *projectService {
	ps := new(projectService)
	ws := new(restful.WebService)

	ws.Path(apiRoot + "/projects").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(ps.listProjects).
		Doc("List projects").
		Param(ws.QueryParameter("from", "minimum identifier of a project")).
		Param(ws.QueryParameter("to", "maximum identifier of a project")).
		Writes(Project{}))

	ws.Route(ws.GET("/{project-id}").To(ps.findProject).
		Doc("Get a project").
		Param(ws.PathParameter("project-id", "identifier of the project")).
		Writes(Project{}))

	ws.Route(ws.POST("").To(ps.createProject).
		Doc("Create a project").
		Reads(Project{}))

	ws.Route(ws.PUT("/{project-id}").To(ps.updateProject).
		Doc("Update a project").
		Param(ws.PathParameter("project-id", "identifier of the project")))

	ws.Route(ws.DELETE("/{project-id}").To(ps.removeProject).
		Doc("Delete a project").
		Param(ws.PathParameter("project-id", "identifier of the project")))

	ps.WebService = ws
	ps.repository = repository

	// set initial stats
	projects, _ := repository.Scan(0, 0)
	stats.ChangeStat("projects", len(projects))

	return ps
}

func (ps *projectService) listProjects(request *restful.Request, response *restful.Response) {
	from, _ := strconv.Atoi(request.QueryParameter("from"))
	to, _ := strconv.Atoi(request.QueryParameter("to"))

	if projects, err := ps.repository.Scan(from, to); err == nil {
		response.WriteEntity(projects)
	} else {
		response.WriteError(http.StatusBadRequest, err)
	}
}

func (ps *projectService) findProject(request *restful.Request, response *restful.Response) {
	id, err := strconv.Atoi(request.PathParameter("project-id"))
	if err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}
	project, err := ps.repository.Get(id)

	if project.Id == 0 {
		response.WriteError(http.StatusNotFound, nil)
	} else {
		response.WriteEntity(project)
	}
}

func (ps *projectService) updateProject(request *restful.Request, response *restful.Response) {
	project := new(Project)
	err := request.ReadEntity(&project)
	if err == nil {
		ps.repository.Put(*project)
		response.WriteEntity(project)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func (ps *projectService) createProject(request *restful.Request, response *restful.Response) {
	id, err := strconv.Atoi(request.PathParameter("project-id"))
	if err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	project := Project{Id: id}
	err = request.ReadEntity(&project)
	if err == nil {
		ps.repository.Put(project)
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(project)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func (ps *projectService) removeProject(request *restful.Request, response *restful.Response) {
	id, err := strconv.Atoi(request.PathParameter("project-id"))
	if err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	err = ps.repository.Remove(id)
	if err == nil {
		response.WriteEntity("removed")
	} else {
		response.WriteError(http.StatusBadRequest, err)
	}
}
