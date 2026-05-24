package middle

import (
	"net/http"
	"dashboard-api/utils"

	"github.com/gin-gonic/gin"
)

//JWTAuth 中间间, 检查token
func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//对登录口放行
		if (len(ctx.Request.URL.String()) >=10 && ctx.Request.URL.String()[0:10] == "/api/login") {
			ctx.Next()
		} else {
			//获取Header中的Authorization
			token := ctx.Request.Header.Get("Authorization")
			if token == "" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "请求未携带token, 无权访问",
					"data": nil,
				})
				ctx.Abort()
				return
			}
			//parseToken解析token包含的信息
			claims, err := utils.JWTToken.ParseToken(token)
			if err != nil {
				//token延期错误
				if err.Error() == "TokenExpired" {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"msg": "收取已过期",
						"data": nil,
					})
					ctx.Abort()
					return
				}
				//其他解析错误
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": err.Error(),
					"data": nil,
				})
				ctx.Abort()
				return
			}
			//继续交由下一个路由处理, 并将解析出的信息传递下去
			ctx.Set("claims", claims)
			ctx.Next()
		}
	}
}