package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	client "github.com/opensource-f2f/open-podcasts/generated/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"sort"
)

type Episode struct {
	pathParam *restful.Parameter
	rssQuery  *restful.Parameter
}

func (r Episode) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/episodes")

	// set the parameters
	r.pathParam = restful.PathParameter("episode", "episode id")
	r.rssQuery = restful.QueryParameter("rss", "The RSS id").Required(true)

	// set the routes
	ws.Route(ws.GET("/").
		Param(r.rssQuery).
		To(r.findAll).
		Returns(http.StatusOK, "OK", []RSS{}))
	ws.Route(ws.GET("/{episode}").
		Param(r.pathParam).
		To(r.findOne).
		Returns(http.StatusOK, "OK", []RSS{}))
	return
}

func (r Episode) findAll(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/rick/.kube/config")
	if err != nil {
		panic(err.Error())
	}
	rss := request.QueryParameter(r.rssQuery.Data().Name)

	ctx := context.Background()
	clientset, err := client.NewForConfig(config)
	spisodeList, err := clientset.Osf2fV1alpha1().Episodes(ns).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("rss=%s", rss),
	})

	sortWithDesc(spisodeList.Items)
	data, err := json.Marshal(spisodeList.Items)
	response.Write(data)
}

func (r Episode) findOne(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/rick/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	name := request.PathParameter(r.pathParam.Data().Name)

	ctx := context.Background()
	clientset, err := client.NewForConfig(config)
	episode, err := clientset.Osf2fV1alpha1().Episodes(ns).Get(ctx, name, metav1.GetOptions{})

	data, err := json.Marshal(episode)
	response.Write(data)
}

func sortWithDesc(items []v1alpha1.Episode) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Spec.Date.After(items[j].Spec.Date.Time)
	})
}
