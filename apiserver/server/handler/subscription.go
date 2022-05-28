package handler

import (
	"context"
	"github.com/emicklei/go-restful/v3"
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/common"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type Subscription struct {
	common.CommonOption
	subscriptionPath *restful.Parameter
}

func (r Subscription) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/subscriptions")

	// set the parameters
	r.subscriptionPath = restful.PathParameter("subscription", "The name of a subscription")

	// set routes
	ws.Route(ws.GET("/{subscription}").
		Param(r.subscriptionPath).
		To(r.findOne).
		Returns(http.StatusOK, "ok", v1alpha1.Subscription{}))
	return
}

func (r Subscription) findOne(req *restful.Request, resp *restful.Response) {
	subName := req.PathParameter(r.subscriptionPath.Data().Name)

	ctx := context.Background()
	subscription, _ := r.Client.Osf2fV1alpha1().Subscriptions(r.DefaultNamespace).Get(ctx, subName, metav1.GetOptions{})
	resp.WriteAsJson(subscription)
}

func uniqueAppend(list []v1.LocalObjectReference, reference v1.LocalObjectReference) (result []v1.LocalObjectReference) {
	found := false
	for i := range list {
		ref := list[i]
		if ref.Name == reference.Name {
			found = true
			break
		}
	}
	result = list
	if !found {
		result = append(list, reference)
	}
	return
}
