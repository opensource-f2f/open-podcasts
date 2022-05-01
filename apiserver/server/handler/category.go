package handler

import (
	"context"
	"encoding/json"
	"github.com/emicklei/go-restful/v3"
	client "github.com/opensource-f2f/open-podcasts/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
)

type Category struct {
}

func (r Category) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/categories")

	ws.Route(ws.GET("/").
		To(r.findAll).
		Returns(http.StatusOK, "OK", []RSS{}))
	return
}

func (r Category) findAll(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/rick/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	clientset, err := client.NewForConfig(config)
	categoryList, err := clientset.MyV1alpha1().Categories(ns).List(ctx, metav1.ListOptions{})

	data, err := json.Marshal(categoryList)
	response.Write(data)
}
