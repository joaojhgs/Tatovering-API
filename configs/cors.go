package configs

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        method := ctx.Request.Method
        if method == "OPTIONS" {
            ctx.Header("Access-Control-Max-Age", "1728000")
            ctx.Header("Access-Control-Allow-Credentials", "true")
            ctx.Header("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
            ctx.Header("Access-Control-Allow-Origin", ctx.Request.Header.Get("origin"))
            ctx.Header("Access-Control-Allow-Headers", "Content-Type,Cookie,Authorization,Access-Control-Request-Headers,Access-Control-Request-Method,Origin,Referer,Sec-Fetch-Dest,Accept-Language,Accept-Encoding,Sec-Fetch-Mode,Sec-Fetch-Site,User-Agent,Pragma,Host,Connection,Cache-Control,Accept-Language,Accept-Encoding,X-Requested-With,X-Forwarded-For,X-Forwarded-Host,X-Forwarded-Proto,X-Forwarded-Port,X-Forwarded-Prefix,X-Real-IP,Accept")
            ctx.AbortWithStatus(http.StatusNoContent)
            return
        }
        ctx.Header("Access-Control-Allow-Origin", ctx.Request.Header.Get("origin"))
        ctx.Header("Access-Control-Allow-Credentials", "true")
        ctx.Next()
    }
}