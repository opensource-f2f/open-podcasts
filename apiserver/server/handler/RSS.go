package handler

import (
	"context"
	"encoding/json"
	"github.com/emicklei/go-restful/v3"
	_ "github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	client "github.com/opensource-f2f/open-podcasts/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
)

const ns = "default"

type RSS struct {
	pathParam          *restful.Parameter
	queryCategoryParam *restful.Parameter
}

func (r RSS) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/rsses")

	r.pathParam = restful.PathParameter("rss", "rss id")
	r.queryCategoryParam = restful.QueryParameter("category", "The category of RSSes")

	ws.Route(ws.GET("/").
		Param(r.queryCategoryParam).
		To(r.findAll).
		Returns(http.StatusOK, "OK", []RSS{}))
	ws.Route(ws.GET("/{rss}").
		Param(r.pathParam).
		To(r.findOne).
		Returns(http.StatusOK, "OK", []RSS{}))
	return
}

func (r RSS) findAll(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/rick/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	clientset, err := client.NewForConfig(config)
	rssList, err := clientset.MyV1alpha1().RSSes(ns).List(ctx, metav1.ListOptions{})

	var filter rssFilter
	if categoryQuery := request.QueryParameter(r.queryCategoryParam.Data().Name); categoryQuery != "" {
		filter = &rssCategoryFilter{keyword: categoryQuery}
	} else {
		filter = &rssNonFilter{}
	}

	data, err := json.Marshal(filter.filter(rssList.Items))
	response.Write(data)
}

func (r RSS) findOne(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/rick/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	name := request.PathParameter(r.pathParam.Data().Name)

	ctx := context.Background()
	clientset, err := client.NewForConfig(config)
	rss, err := clientset.MyV1alpha1().RSSes(ns).Get(ctx, name, metav1.GetOptions{})

	data, err := json.Marshal(rss)
	response.Write(data)
}
