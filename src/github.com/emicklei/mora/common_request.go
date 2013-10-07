package main

import (
	"github.com/emicklei/go-restful"

	"log"
	"net/http"
)

func handleError(err error, resp *restful.Response) {
	if err.Error() == "not found" {
		resp.WriteError(http.StatusNotFound, err)
		return
	}
	if err.Error() == "unauthorized" {
		resp.WriteError(http.StatusUnauthorized, err)
		return
	}
	log.Printf("[mora] error:%v", err)
	resp.AddHeader("Content-Type", "text/plain") // consider making ServiceError and write JSON
	resp.WriteErrorString(500, err.Error())
}

func optionsOK(req *restful.Request, resp *restful.Response) {
	resp.WriteHeader(http.StatusOK)
}
