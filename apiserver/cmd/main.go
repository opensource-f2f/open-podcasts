package main

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/common"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/filter"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/handler"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/user"
	client "github.com/opensource-f2f/open-podcasts/generated/clientset/versioned"
	"github.com/spf13/cobra"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
)

func main() {
	cmd := newCommand()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

type option struct {
	defaultNamespace string
}

func newCommand() (cmd *cobra.Command) {
	opt := &option{}
	cmd = &cobra.Command{
		Use:  "apiserver",
		RunE: opt.runE,
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.defaultNamespace, "default-namespace", "", "default",
		"The default namespace of the resources")
	return
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	var config *restclient.Config
	config, err = clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		return
	}

	var clientset *client.Clientset
	if clientset, err = client.NewForConfig(config); err != nil {
		return
	}
	commonOption := common.CommonOption{
		Client:           clientset,
		DefaultNamespace: o.defaultNamespace,
	}

	rss := &handler.RSS{CommonOption: commonOption}
	category := &handler.Category{CommonOption: commonOption}
	episode := &handler.Episode{CommonOption: commonOption}
	profile := &handler.Profile{CommonOption: commonOption}
	subscription := &handler.Subscription{CommonOption: commonOption}
	notifier := &handler.Notifier{CommonOption: commonOption}
	showItem := &handler.ShowItem{CommonOption: commonOption}
	userws := &user.User{}

	restful.DefaultContainer.Add(rss.WebService())
	restful.DefaultContainer.Add(category.WebService())
	restful.DefaultContainer.Add(episode.WebService())
	restful.DefaultContainer.Add(profile.WebService())
	restful.DefaultContainer.Add(subscription.WebService())
	restful.DefaultContainer.Add(notifier.WebService())
	restful.DefaultContainer.Add(showItem.WebService())
	restful.DefaultContainer.Add(userws.WebService())

	setupFilter(restful.RegisteredWebServices())

	restConfig := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		APIPath:     "/apidocs.json",
	}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(restConfig))

	log.Printf("start listening on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
	return
}

func setupFilter(wss []*restful.WebService) {
	for i := range wss {
		ws := wss[i]
		ws.Filter(filter.AuthJWT)
	}
}
