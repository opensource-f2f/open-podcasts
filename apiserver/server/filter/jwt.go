package filter

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

var (
	sharedSecret = []byte("shared-token")
)

func AuthJWT(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	authHeader := req.HeaderParameter("Authorization")

	if !validJWT(authHeader) {
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}

	chain.ProcessFilter(req, resp)
}

func validJWT(authHeader string) bool {
	if !strings.HasPrefix(authHeader, "JWT ") {
		return false
	}

	jwtToken := strings.Split(authHeader, " ")
	if len(jwtToken) < 2 {
		return false
	}
	parts := strings.Split(jwtToken[1], ".")
	err := jwt.SigningMethodHS512.Verify(strings.Join(parts[0:2], "."), parts[2], sharedSecret)
	if err != nil {
		return false
	}

	return true
}
