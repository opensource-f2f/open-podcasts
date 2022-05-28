package handler

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type Notifier struct {
	common.CommonOption
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
	name := request.PathParameter(r.notifierPath.Data().Name)

	ctx := context.Background()
	notifier, err := r.Client.Osf2fV1alpha1().Notifiers(r.DefaultNamespace).Get(ctx, name, metav1.GetOptions{})
	fmt.Println(err)
	response.WriteAsJson(notifier)
}
