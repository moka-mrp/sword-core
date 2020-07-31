package middleware

import (
   "github.com/gin-gonic/gin"
)

// CrossDomain 全局添加跨域允许
func Cors() gin.HandlerFunc {
   return func(c *gin.Context) {
      c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
      c.Writer.Header().Set("Access-Control-Allow-Headers", "access-token, x-requested-with, content-type")
      c.Writer.Header().Set("Access-Control-Allow-Methods","POST, GET, OPTIONS, PUT, PATCH, DELETE")
   }
}