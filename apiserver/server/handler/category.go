package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type Category struct {
	common.CommonOption
}

func (r Category) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/categories")

	ws.Route(ws.GET("/").
		To(r.findAll).
		Returns(http.StatusOK, "OK", []RSS{}))
	return
}

func (r *Category) findAll(request *restful.Request, response *restful.Response) {
	ctx := context.Background()
	categoryList, err := r.Client.Osf2fV1alpha1().Categories(r.DefaultNamespace).List(ctx, metav1.ListOptions{})
	fmt.Println(err)

	data, err := json.Marshal(categoryList)
	response.Write(data)
}
