package main

import (
	"flag"
	"github.com/dmotylev/goproperties"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/emicklei/landskape/application"
	"github.com/emicklei/landskape/dao"
	"github.com/emicklei/landskape/webservice"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

var propertiesFile = flag.String("config", "landskape.properties", "the configuration file")

func main() {
	log.Print("[landskape] service startup...")
	flag.Parse()
	props, _ := properties.Load(*propertiesFile)
	session, _ := mgo.Dial(props["mongo.connection"]) // TODO error checking
	defer session.Close()

	appDao := dao.SystemDao{session.DB(props["mongo.database"]).C("systems")}
	conDao := dao.ConnectionDao{session.DB(props["mongo.database"]).C("connections")}
	application.SharedLogic = application.Logic{appDao, conDao}

	webservice.SystemResource{application.SharedLogic}.Register()
	webservice.ConnectionResource{application.SharedLogic}.Register()

	// graphical diagrams
	restful.Add(webservice.NewDiagramService())
	webservice.DotConfig["binpath"] = props["dot.path"]
	webservice.DotConfig["tmp"] = props["dot.tmp"]

	// expose api using swagger
	basePath := "http://" + props["http.server.host"] + ":" + props["http.server.port"]
	config := swagger.Config{
		WebServicesUrl:  basePath,
		ApiPath:         props["swagger.api"],
		SwaggerPath:     props["swagger.path"],
		SwaggerFilePath: props["swagger.home"],
		WebServices:     restful.RegisteredWebServices()}
	swagger.InstallSwaggerService(config)

	log.Printf("[landskape] ready to serve on %v\n", basePath)
	log.Fatal(http.ListenAndServe(":"+props["http.server.port"], nil))
}
