package main

import (
	"encoding/json"
	"fmt"
	"github.com/ku-ovdp/api/endpoints"
	"github.com/ku-ovdp/api/persistence"
	_ "github.com/ku-ovdp/api/persistence/dummy"
	_ "github.com/ku-ovdp/api/persistence/mgo"
	"github.com/ku-ovdp/api/repository"
	"github.com/ku-ovdp/api/sockgroup"
	"github.com/ku-ovdp/api/stats"
	"github.com/traviscline/go-restful"
	"log"
	"net/http"
	"time"
)

// Create application services and dependancies
func constructApplication() {

	// construct backend
	backend := persistence.Get(*persistenceBackend)
	if backend == nil {
		log.Fatalln("Invalid repository backend.", *persistenceBackend)
	}
	backend.Init()

	// build repositories
	repositories := repository.NewRepositoryGroup()
	projectRepository := backend.NewProjectRepository(repositories)
	repositories["projects"] = projectRepository

	sessionRepository := backend.NewSessionRepository(repositories)
	repositories["sessions"] = sessionRepository

	sampleRepository := backend.NewSampleRepository(repositories)
	repositories["samples"] = sampleRepository

	// register logging Dispatch method
	restful.Dispatch = loggingDispatch
	// set default mime type
	restful.DefaultResponseMimeType = restful.MIME_JSON

	// construct and register services
	apiRoot := fmt.Sprintf("/v%d", API_VERSION)
	restful.Add(endpoints.NewProjectService(apiRoot, projectRepository))
	restful.Add(endpoints.NewSessionService(apiRoot, sessionRepository))
	restful.Add(endpoints.NewVoiceSampleService(apiRoot, sampleRepository))

	// sockgroup for publishing statistics
	sg := sockgroup.NewGroup()
	sg.Start()
	stats.Destination = sg
	http.Handle("/v1/livestats/", sg.Handler("/v1/livestats"))

	// other handlers
	http.HandleFunc("/v1/stats/", stats.Handler)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", indexHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(struct {
		Url string `json:"documentation_url"`
	}{
		fmt.Sprintf("/v%d-docs", API_VERSION),
	})
	fmt.Fprintln(w, string(b))
}

func loggingDispatch(w http.ResponseWriter, r *http.Request) {
	lwr := &loggedResponseWriter{w, 0}
	t1 := time.Now()
	restful.DefaultDispatch(lwr, r)
	fmt.Println(r.Method, r.URL, lwr.status, time.Now().Sub(t1))
}

type loggedResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *loggedResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
	stats.ChangeStat("requests", 1)
}

func (w *loggedResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
}
