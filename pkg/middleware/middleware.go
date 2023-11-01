package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	middlewares []gin.HandlerFunc
}

func (m *Middleware) Get() []gin.HandlerFunc {
	return m.middlewares
}

func (m *Middleware) AllowConfig(allowAddress, allowMethods, allowHeaders []string, allowCredentials bool) gin.HandlerFunc {
	cfg := cors.DefaultConfig()
	cfg.AllowOrigins = allowAddress
	cfg.AllowMethods = allowMethods
	cfg.AllowHeaders = allowHeaders
	cfg.AllowCredentials = allowCredentials
	return cors.New(cfg)
}