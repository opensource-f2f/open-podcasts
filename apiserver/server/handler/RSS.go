package handler

import (
	"bytes"
	"context"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	_ "github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	"github.com/opensource-f2f/open-podcasts/apiserver/server/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"sigs.k8s.io/yaml"
	"text/template"
	"time"
)

type RSS struct {
	common.CommonOption
	pathParam          *restful.Parameter
	queryCategoryParam *restful.Parameter
}

func (r RSS) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/rsses")

	r.pathParam = restful.PathParameter("rss", "rss id")
	r.queryCategoryParam = restful.QueryParameter("category", "The category of RSSes")

	ws.Route(ws.POST("/").
		Param(r.queryCategoryParam).
		To(r.create).
		Returns(http.StatusOK, "OK", []RSS{}))
	ws.Route(ws.GET("/").
		Param(r.queryCategoryParam).
		To(r.findAll).
		Returns(http.StatusOK, "OK", []RSS{}))
	ws.Route(ws.GET("/{rss}").
		Param(r.pathParam).
		To(r.findOne).
		Returns(http.StatusOK, "OK", []RSS{}))
	ws.Route(ws.GET("/opml/export").
		To(r.opmlExport).
		Returns(http.StatusOK, "OK", []RSS{}))
	ws.Route(ws.GET("/yaml/export").
		To(r.yamlExport).
		Returns(http.StatusOK, "OK", []RSS{}))
	return
}

func (r RSS) create(request *restful.Request, response *restful.Response) {
	ctx := context.Background()

	rss := &v1alpha1.RSS{}
	request.ReadEntity(rss) // TODO handle the error case

	rss.GenerateName = "auto"
	r.Client.Osf2fV1alpha1().RSSes(r.DefaultNamespace).Create(ctx, rss, metav1.CreateOptions{})
	response.Write([]byte("ok"))
}

func (r RSS) findAll(request *restful.Request, response *restful.Response) {
	ctx := context.Background()
	rssList, err := r.Client.Osf2fV1alpha1().RSSes(r.DefaultNamespace).List(ctx, metav1.ListOptions{})
	fmt.Println(err)

	var filter rssFilter
	if categoryQuery := request.QueryParameter(r.queryCategoryParam.Data().Name); categoryQuery != "" {
		filter = &rssCategoryFilter{keyword: categoryQuery}
	} else {
		filter = &rssNonFilter{}
	}
	response.WriteAsJson(filter.filter(rssList.Items))
}

func (r RSS) findOne(request *restful.Request, response *restful.Response) {
	name := request.PathParameter(r.pathParam.Data().Name)

	ctx := context.Background()
	rss, err := r.Client.Osf2fV1alpha1().RSSes(r.DefaultNamespace).Get(ctx, name, metav1.GetOptions{})
	fmt.Println(err)
	response.WriteAsJson(rss)
}

func (r RSS) opmlExport(request *restful.Request, response *restful.Response) {
	ctx := context.Background()
	rssList, err := r.Client.Osf2fV1alpha1().RSSes(r.DefaultNamespace).List(ctx, metav1.ListOptions{})
	response.Header().Set("Content-Type", "application/octet-stream")
	response.Header().Set("Content-Disposition", "attachment; filename=rsses.opml")

	fmt.Println(err)
	data, err := (&opmlRSSList{
		Title:      "Open Podcast",
		CreateDate: time.Now(),
		Items:      rssList.Items,
	}).asXML()
	response.Write(data)
}

type opmlRSSList struct {
	Title      string
	CreateDate time.Time
	Items      []v1alpha1.RSS
}

func (o *opmlRSSList) asXML() (xml []byte, err error) {
	var tpl *template.Template
	if tpl, err = template.New("opml").Parse(opmlTemplate); err == nil {
		buf := bytes.NewBuffer([]byte{})
		if err = tpl.Execute(buf, o); err == nil {
			xml = buf.Bytes()
		}
	}
	return
}

var opmlTemplate = `<?xml version='1.0' encoding='UTF-8' standalone='no' ?>
<opml version="2.0">
  <head>
    <title>{{.Title}}</title>
    <dateCreated>{{.CreateDate}}</dateCreated>
  </head>
  <body>
	{{- range $val := .Items}}
    <outline text="{{$val.Spec.Description}}" title="{{$val.Spec.Title}}" type="rss" xmlUrl="{{$val.Spec.Address}}" htmlUrl="{{$val.Spec.Link}}"/>
	{{- end}}
  </body>
</opml>
`

func (r RSS) yamlExport(request *restful.Request, response *restful.Response) {
	ctx := context.Background()
	rssList, err := r.Client.Osf2fV1alpha1().RSSes(r.DefaultNamespace).List(ctx, metav1.ListOptions{})
	fmt.Println(err)
	response.Header().Set("Content-Type", "application/octet-stream")
	response.Header().Set("Content-Disposition", "attachment; filename=rsses.yaml")
	data, err := yaml.Marshal(rssList)
	response.Write(data)
}
