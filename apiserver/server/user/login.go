package user

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

var (
	sharedSecret = []byte("shared-token")
)

type User struct {
}

func (u User) WebService() (ws *restful.WebService) {
	ws = new(restful.WebService)
	ws.Path("/auth")

	ws.Route(ws.GET("/login").
		To(u.login).
		Returns(http.StatusOK, "OK", []User{}))
	return
}

func (u User) login(request *restful.Request, response *restful.Response) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user": "xxx",
	})

	tokenString, err := token.SignedString(sharedSecret)

	fmt.Println(err)
	response.Write([]byte(tokenString))
}
