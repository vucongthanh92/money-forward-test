package main

import (
	"os"

	"money-forward-test/routers"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func getPort() string {
	p := os.Getenv("HOST_PORT")
	if p != "" {
		return ":" + p
	}
	return ":3040"
}

// CORSMiddleware func
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Cache-Control, X-Requested-With, content-type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(location.Default())
	r.Use(CORSMiddleware())
	rg := r.Group("/api")
	rg.Use(CORSMiddleware())
	{
		routers.Route(rg)
	}
	return r
}

func main() {
	port := getPort()
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := setupRouter()
	r.Run(port)
}
