package handler

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	client "github.com/opensource-f2f/open-podcasts/generated/clientset/versioned"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

type ShowItem struct {
	showitemPath *restful.Parameter
	showQuery    *restful.Parameter
}

func (r ShowItem) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/showitems")

	// set the parameters
	r.showitemPath = restful.PathParameter("showitem", "")
	r.showQuery = restful.QueryParameter("show", "The name of a show")

	// set the routes
	ws.Route(ws.POST("/").
		To(r.create).
		Reads(&v1alpha1.ShowItem{}).
		Returns(http.StatusOK, "OK", []RSS{}))
	ws.Route(ws.GET("/").
		To(r.findAll).
		Param(r.showQuery).
		Returns(http.StatusOK, "OK", []v1alpha1.ShowItem{}))
	ws.Route(ws.POST("/{showitem}/upload").
		To(r.upload).
		Consumes("multipart/form-data").
		Returns(http.StatusOK, "OK", []RSS{}))
	ws.Route(ws.GET("/{showitem}/download").
		To(r.download).
		Returns(http.StatusOK, "OK", []RSS{}))
	return
}

func (r ShowItem) create(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	clientset, err := client.NewForConfig(config)

	showitem := &v1alpha1.ShowItem{}
	request.ReadEntity(showitem)

	if len(showitem.Labels) == 0 {
		showitem.Labels = map[string]string{}
	}
	showitem.Labels[v1alpha1.LabelKeyShowRef] = showitem.Spec.ShowRef

	// calculate the next index
	var showItemList *v1alpha1.ShowItemList
	if showItemList, err = clientset.Osf2fV1alpha1().ShowItems(ns).List(ctx, v1.ListOptions{
		LabelSelector: fmt.Sprintf("rss=%s", showitem.Spec.ShowRef),
	}); err == nil {
		showitem.Spec.Index = len(showItemList.Items)
	}

	_, _ = clientset.Osf2fV1alpha1().ShowItems(ns).Create(ctx, showitem, v1.CreateOptions{})
}

func (r ShowItem) findAll(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	clientset, err := client.NewForConfig(config)

	showName := request.QueryParameter(r.showQuery.Data().Name)
	labelSelector := ""
	if showName != "" {
		labelSelector = fmt.Sprintf("%s=%s", v1alpha1.LabelKeyShowRef, showName)
	}

	showItemList, _ := clientset.Osf2fV1alpha1().ShowItems(ns).List(ctx, v1.ListOptions{
		LabelSelector: labelSelector,
	})
	response.WriteAsJson(showItemList)
}

func (r ShowItem) upload(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}

	showItemName := request.PathParameter(r.showitemPath.Data().Name)
	ctx := context.Background()
	clientset, err := client.NewForConfig(config)
	showItem, err := clientset.Osf2fV1alpha1().ShowItems(ns).Get(ctx, showItemName, v1.GetOptions{})

	httpReq := request.Request
	err = httpReq.ParseMultipartForm(10 << 20)

	for fName := range httpReq.MultipartForm.File {
		var f multipart.File
		f, _, err = httpReq.FormFile(fName)
		defer func() {
			_ = f.Close()
		}()

		var data []byte
		data, err = ioutil.ReadAll(f)

		filepath := path.Join(getTempDir(), showItem.Spec.Filename)
		ioutil.WriteFile(filepath, data, 0644)
		break
	}
	fmt.Println(err)
}

func (r ShowItem) download(request *restful.Request, response *restful.Response) {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}

	showItemName := request.PathParameter(r.showitemPath.Data().Name)

	ctx := context.Background()
	clientset, err := client.NewForConfig(config)
	showItem, err := clientset.Osf2fV1alpha1().ShowItems(ns).Get(ctx, showItemName, v1.GetOptions{})
	if err == nil {
		filepath := path.Join(getTempDir(), showItem.Spec.Filename)
		response.Header().Set("Content-Type", "application/octet-stream")
		response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", showItem.Spec.Filename))

		data, err := ioutil.ReadFile(filepath)
		if err == nil {
			_, _ = response.Write(data)
		} else {
			_, _ = response.Write([]byte(err.Error()))
		}
	}
}

func getTempDir() string {
	return os.TempDir()
}
