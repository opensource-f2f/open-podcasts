package handler

import (
	"context"
	"github.com/emicklei/go-restful/v3"
	client "github.com/opensource-f2f/open-podcasts/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
)

type Notifier struct {
	notifierPath *restful.Parameter
}

func (r Notifier) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/notifiers")

	// set the parameters
	r.notifierPath = restful.PathParameter("notifier", "The name of a notifier")

	// set the routes
	ws.Route(ws.GET("/{notifier}").
		Param(r.notifierPath).
		To(r.findOne).
		Returns(http.StatusOK, "OK", []RSS{}))
	return
}

func (r Notifier) findOne(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/rick/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	name := request.PathParameter(r.notifierPath.Data().Name)

	ctx := context.Background()
	clientset, err := client.NewForConfig(config)
	notifier, err := clientset.MyV1alpha1().Notifiers(ns).Get(ctx, name, metav1.GetOptions{})
	response.WriteAsJson(notifier)
}
