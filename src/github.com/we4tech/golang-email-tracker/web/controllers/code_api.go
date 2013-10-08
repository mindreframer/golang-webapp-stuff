package controllers

import (
	"net/http"
	"encoding/json"

	"appengine"
	"models"
	"strconv"

	"github.com/gorilla/sessions"
	"io"
	"bytes"
)

func (c *AppController) ApiCode(w http.ResponseWriter, r *http.Request, session *sessions.Session) *ActionResponse {
	ctx := appengine.NewContext(r)
	user := c.GetCurrentUser(session)

	switch r.Method {
	case "POST":
		if r.FormValue("track") == "" {
			return createCode(ctx, user, r)
		} else {
			return trackCode(ctx, user, r)
		}
	case "GET":
		return retrieveAllCodes(ctx, user, r)
	case "DELETE":
		return destroyCode(ctx, user, r)
	default:
		return &ActionResponse{
			RenderText: "Hello there",
		}
	}
}

func destroyCode(ctx appengine.Context, user *models.User, r *http.Request) *ActionResponse {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	removed, err := models.DestroyCodeBy(ctx, id)

	if removed {
		return (&ActionResponse{
				Code: http.StatusNoContent,
				RenderText: "",
			}).AsJson()
	} else {
		return (&ActionResponse{
				Code: http.StatusNotFound,
				RenderText: createApiError(err).AsJson(),
			}).AsJson()
	}
}

func trackCode(ctx appengine.Context, user *models.User, r *http.Request) *ActionResponse {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	updated, err := models.StartTrackingCode(ctx, id)

	if updated {
		return (&ActionResponse{
				Code: http.StatusNoContent,
				RenderText: "",
			}).AsJson()
	} else {
		return (&ActionResponse{
				Code: http.StatusNotFound,
				RenderText: createApiError(err).AsJson(),
			}).AsJson()
	}
}

func retrieveAllCodes(ctx appengine.Context, user *models.User, r *http.Request) *ActionResponse {
	codes, err := models.FindCodesByUserId(ctx, user.Id)
	content := ""

	if err != nil {
		content = createApiError(err).AsJson()
	} else {
		d, _ := json.Marshal(codes)
		content = string(d)
	}

	return (&ActionResponse{
			Code: http.StatusOK,
			RenderText: content,
		}).AsJson()
}

func createCode(ctx appengine.Context, user *models.User, r *http.Request) *ActionResponse {
	code := initNewCode(r, user)
	saved, validationErr, err := code.Save(ctx)
	httpCode := http.StatusBadRequest
	content := ""

	if saved {

		content = code.AsJson()
		httpCode = http.StatusCreated

	} else if validationErr != nil {

		content = createApiError(validationErr).AsJson()

	} else if err != nil {

		content = createApiError(err).AsJson()
	}

	return (&ActionResponse{
			Code: httpCode,
			RenderText: content,
		}).AsJson()
}

type JsonFields struct {
	Key,  Value string
}

func initNewCode(r *http.Request, user *models.User) *models.Code {
	data := readBody(r)
	c := new(models.Code)
	c.Title = data["title"].(string)
	c.UserId = user.Id

	return c
}

func readBody(r *http.Request) map[string]interface {} {
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, r.Body)

	data := make(map[string]interface{})
	json.Unmarshal(buf.Bytes(), &data)

	return data
}
