package handler

import (
	"github.com/emicklei/go-restful/v3"
	"net/http"
)

type Simple interface {
	findOne(request *restful.Request, response *restful.Response)
}

func output(data interface{}, err error, response *restful.Response) {
	if err != nil {
		_ = response.WriteError(http.StatusBadGateway, err)
	} else {
		_ = response.WriteAsJson(data)
	}
}
