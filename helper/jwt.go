package helper

import (
	"errors"
	"ta-kasir/base"
	"ta-kasir/model/request"

	"github.com/gin-gonic/gin"
)

func GetClaims(c *gin.Context) (*request.JwtClaim, error) {
	claims, ok := c.Get("jwt_claims")
	if !ok {
		return nil, errors.New(base.NoClaimsFound)
	}
	typeAssert, ok := claims.(request.JwtClaim)
	if !ok {
		return nil, errors.New(base.NoClaimsFound)
	}
	return &typeAssert, nil
}