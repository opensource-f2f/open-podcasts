package handler

import (
	"context"
	"github.com/emicklei/go-restful/v3"
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	client "github.com/opensource-f2f/open-podcasts/generated/clientset/versioned"
	"io/ioutil"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
)

type ShowItem struct {
}

func (r ShowItem) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/showitems")

	// set the parameters

	// set the routes
	ws.Route(ws.POST("/").
		To(r.create).
		Reads(&v1alpha1.ShowItem{}, "").
		Returns(http.StatusOK, "OK", []RSS{}))
	ws.Route(ws.POST("/{showitem}/upload").
		To(r.upload).
		Consumes("multipart/form-data").
		Returns(http.StatusOK, "OK", []RSS{}))
}

func (r ShowItem) create(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/rick/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	clientset, err := client.NewForConfig(config)
	_, _ = clientset.MyV1alpha1().Profiles(ns).Create(ctx, &v1alpha1.Profile{}
}

func (r ShowItem) upload(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/rick/.kube/config")
	if err != nil {
		panic(err.Error())
	}

	httpReq := request.Request
	err := httpReq.ParseMultipartForm(10 << 20)

	for fName := range httpReq.MultipartForm.File {
		f, _, err := httpReq.FormFile(fName)
		data, err := ioutil.ReadAll(f)
		f.Close()
		break
	}

	ctx := context.Background()
	clientset, err := client.NewForConfig(config)
	_, _ = clientset.MyV1alpha1().Profiles(ns).Create(ctx, &v1alpha1.Profile{}
}
