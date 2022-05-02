package main

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/filter"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/handler"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/user"
	"log"
	"net/http"
)

func main() {
	rss := &handler.RSS{}
	category := &handler.Category{}
	episode := &handler.Episode{}
	profile := &handler.Profile{}
	subscription := &handler.Subscription{}
	notifier := &handler.Notifier{}
	userws := &user.User{}
	restful.DefaultContainer.Add(rss.WebService())
	restful.DefaultContainer.Add(category.WebService())
	restful.DefaultContainer.Add(episode.WebService())
	restful.DefaultContainer.Add(profile.WebService())
	restful.DefaultContainer.Add(subscription.WebService())
	restful.DefaultContainer.Add(notifier.WebService())
	restful.DefaultContainer.Add(userws.WebService())

	setupFilter(restful.RegisteredWebServices())

	config := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		APIPath:     "/apidocs.json",
	}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	log.Printf("start listening on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func setupFilter(wss []*restful.WebService) {
	for i := range wss {
		ws := wss[i]
		ws.Filter(filter.AuthJWT)
	}
}
