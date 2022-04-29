package main

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/filter"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/handler"
	"log"
	"net/http"
)

func main() {
	rss := &handler.RSS{}
	restful.DefaultContainer.Add(rss.WebService())

	setupFilter(restful.RegisteredWebServices())

	config := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		APIPath:     "/apidocs.json",
	}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	log.Printf("start listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupFilter(wss []*restful.WebService) {
	for i := range wss {
		ws := wss[i]
		ws.Filter(filter.AuthJWT)
	}
}
