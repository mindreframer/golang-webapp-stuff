package endpoints

import (
	"fmt"
	"github.com/kr/s3/s3util"
	. "github.com/ku-ovdp/api/entities"
	. "github.com/ku-ovdp/api/repository"
	"github.com/ku-ovdp/api/stats"
	"github.com/traviscline/go-restful"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type sampleService struct {
	*restful.WebService
	repository VoiceSampleRepository
}

func NewVoiceSampleService(apiRoot string, repository VoiceSampleRepository) *sampleService {
	s := new(sampleService)
	ws := new(restful.WebService)

	ws.Path(apiRoot + "/session/{session-id}/samples").
		Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(s.listVoiceSamples).
		Doc("List voice samples").
		Param(ws.PathParameter("session-id", "identifier of the session").DataType("int")).
		Param(ws.QueryParameter("from", "minimum identifier of a project")).
		Param(ws.QueryParameter("to", "maximum identifier of a project")).
		Writes(VoiceSample{}))

	ws.Route(ws.POST("").To(s.createVoiceSample).
		Doc("Create a voice sample").
		Param(ws.PathParameter("session-id", "identifier of the session").DataType("int")).
		Reads(VoiceSample{}))

	ws.Route(ws.GET("/{sample-index}").To(s.findVoiceSample).
		Doc("Get a voice sample").
		Param(ws.PathParameter("session-id", "identifier of the session").DataType("int")).
		Param(ws.PathParameter("sample-index", "identifier of the sample").DataType("int")).
		Writes(VoiceSample{}))

	ws.Route(ws.PUT("/{sample-index}").To(s.updateVoiceSample).
		Doc("Update a voice sample").
		Param(ws.PathParameter("session-id", "identifier of the session").DataType("int")).
		Param(ws.PathParameter("sample-index", "identifier of the sample").DataType("int")).
		Param(ws.BodyParameter("VoiceSample", "the session entity").DataType("string")))

	ws.Route(ws.DELETE("/{sample-index}").To(s.removeVoiceSample).
		Doc("Delete a voice sample").
		Param(ws.PathParameter("session-id", "identifier of the session").DataType("int")).
		Param(ws.PathParameter("sample-index", "identifier of the sample").DataType("int")))

	ws.Route(ws.GET("/{sample-index}/audio").To(s.streamVoiceSample).
		Doc("Get a voice sample's audio").
		Param(ws.PathParameter("session-id", "identifier of the session").DataType("int")).
		Param(ws.PathParameter("sample-index", "identifier of the sample").DataType("int")).
		Writes(VoiceSample{}))

	ws.Route(ws.PUT("/{sample-index}/audio").To(s.uploadVoiceSample).
		Doc("Attach audio to a voice sample").
		Consumes("*/*").
		Param(ws.PathParameter("session-id", "identifier of the session").DataType("int")).
		Param(ws.PathParameter("sample-index", "identifier of the sample").DataType("int")).
		Param(ws.BodyParameter("Audio", "the audio blob entity").DataType("string")))

	s.WebService = ws
	s.repository = repository

	return s
}

func (s *sampleService) listVoiceSamples(request *restful.Request, response *restful.Response) {
	sessionId, _ := strconv.Atoi(request.PathParameter("session-id"))

	from, _ := strconv.Atoi(request.QueryParameter("from"))
	to, _ := strconv.Atoi(request.QueryParameter("to"))

	if samples, err := s.repository.Scan(sessionId, from, to); err == nil {
		response.WriteEntity(samples)
	} else {
		response.WriteError(http.StatusBadRequest, err)
	}
}

func (s *sampleService) findVoiceSample(request *restful.Request, response *restful.Response) {
	sessionId, _ := strconv.Atoi(request.PathParameter("session-id"))
	sampleId, _ := strconv.Atoi(request.PathParameter("sample-index"))

	sample, _ := s.repository.Get(sessionId, sampleId)

	if sample.Id == 0 {
		response.WriteError(http.StatusNotFound, nil)
	} else {
		response.WriteEntity(sample)
	}
}

func (s *sampleService) createVoiceSample(request *restful.Request, response *restful.Response) {
	sessionId, _ := strconv.Atoi(request.PathParameter("session-id"))

	sample := VoiceSample{
		SessionId: sessionId,
		Created:   time.Now(),
	}
	var err error
	sample, err = s.repository.Put(sample)
	if err != nil {
		response.WriteError(http.StatusBadRequest, err)
		return
	}

	stats.ChangeStat("samples", 1)
	response.WriteHeader(http.StatusCreated)
	response.WriteEntity(sample)
}

func (s *sampleService) updateVoiceSample(request *restful.Request, response *restful.Response) {
	sample := new(VoiceSample)
	err := request.ReadEntity(&sample)
	if err == nil {
		s.repository.Put(*sample)
		response.WriteEntity(sample)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func (s *sampleService) removeVoiceSample(request *restful.Request, response *restful.Response) {
	sessionId, _ := strconv.Atoi(request.PathParameter("session-id"))
	sampleId, _ := strconv.Atoi(request.PathParameter("sample-index"))

	err := s.repository.Remove(sessionId, sampleId)
	if err == nil {
		response.WriteEntity("removed")
	} else {
		response.WriteError(http.StatusBadRequest, err)
	}
}

func (s *sampleService) streamVoiceSample(request *restful.Request, response *restful.Response) {
	sessionId, _ := strconv.Atoi(request.PathParameter("session-id"))
	sampleId, _ := strconv.Atoi(request.PathParameter("sample-index"))

	sample, _ := s.repository.Get(sessionId, sampleId)

	if sample.Id == 0 {
		response.WriteError(http.StatusNotFound, nil)
		return
	}
	ak, sk := os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY")
	if ak == "" || sk == "" {
		response.WriteError(http.StatusInternalServerError, fmt.Errorf("Missing S3_ACCESS_KEY and/or S3_SECRET_KEY"))
		return
	}
	s3util.DefaultConfig.AccessKey = ak
	s3util.DefaultConfig.SecretKey = sk

	audio, err := s3util.Open(sample.AudioURL, nil)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	defer audio.Close()
	response.AddHeader("Content-Type", sample.MimeType)
	io.Copy(response, audio)
}

func (s *sampleService) uploadVoiceSample(request *restful.Request, response *restful.Response) {
	sessionId, _ := strconv.Atoi(request.PathParameter("session-id"))
	sampleId, _ := strconv.Atoi(request.PathParameter("sample-index"))

	sample, err := s.repository.Get(sessionId, sampleId)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}
	session, err := s.repository.Group().Sessions().Get(sessionId)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	ak, sk := os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")

	if ak == "" || sk == "" {
		response.WriteError(http.StatusInternalServerError, fmt.Errorf("Missing S3_ACCESS_KEY and/or S3_SECRET_KEY"))
		return
	} else if bucket == "" {
		response.WriteError(http.StatusInternalServerError, fmt.Errorf("No S3_BUCKET env var present."))
		return
	} else if bucket[len(bucket)-1] != '/' {
		response.WriteError(http.StatusInternalServerError, fmt.Errorf("S3_BUCKET must have trailing slash."))
		return
	}
	s3util.DefaultConfig.AccessKey = ak
	s3util.DefaultConfig.SecretKey = sk

	destURI := fmt.Sprintf("%s%d/%d/%d", bucket, session.ProjectId, sessionId, sampleId)
	w, err := s3util.Create(destURI, nil, nil)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	io.Copy(w, request.Request.Body)
	w.Close()

	sample.AudioURL = destURI

	if sample, err = s.repository.Put(sample); err == nil {
		response.WriteEntity(sample)
	} else {
		response.WriteError(http.StatusBadRequest, err)
	}
}
