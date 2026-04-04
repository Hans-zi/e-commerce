package middleware

import (
	"e-commerce/internal/consts"
	"e-commerce/pkg/token"
	"e-commerce/pkg/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(consts.AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("authorization format is invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != consts.AuthorizationHeaderType {
			err := errors.New("authorization type is invalid")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		ctx.Set(consts.AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, ok := ctx.Get(consts.AuthorizationPayloadKey)
		if !ok {
			err := errors.New("payload is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		payload, ok := data.(*token.Payload)
		if !ok {
			err := errors.New("invalid payload format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		
		if payload.Role != consts.ROLE_ADMIN {
			err := errors.New("invalid role, need admin")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		ctx.Next()
	}
}
