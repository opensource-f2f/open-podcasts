package handler

import "github.com/emicklei/go-restful/v3"

type Simple interface {
	findOne(request *restful.Request, response *restful.Response)
}
