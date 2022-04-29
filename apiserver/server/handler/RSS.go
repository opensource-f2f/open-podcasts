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
}

func (r RSS) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/rsses")

	ws.Route(ws.GET("/").
		To(r.findAll).
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
	//scheme.AddToScheme(clientsetscheme.Scheme)
	rssList, err := clientset.MyV1alpha1().RSSs(ns).List(ctx, metav1.ListOptions{})

	data, err := json.Marshal(rssList)
	response.Write(data)
}
