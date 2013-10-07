package main

import (
	"flag"
	"fmt"
	"github.com/traviscline/go-restful"
	"log"
	"path/filepath"
)

var baseUrl = flag.String("baseURL", "http://api.openvoicedata.org", "Base URL")

func installSwagger() {
	flag.Parse()
	swaggerUiPath, err := filepath.Abs("static/swagger-ui/dist/")
	if err != nil {
		log.Fatalln(err)
	}
	config := restful.SwaggerConfig{
		WebServicesUrl:  *baseUrl,
		ApiPath:         fmt.Sprintf("/v%d-docs", API_VERSION),
		SwaggerPath:     "/docs/",
		SwaggerFilePath: swaggerUiPath}
	restful.InstallSwaggerService(config)
}

func init() {
	installSwagger()
}
