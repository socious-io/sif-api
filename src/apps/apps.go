package apps

import (
	"context"
	"fmt"
	"net/http"
	"sif/src/apps/utils"
	"sif/src/apps/views"
	"sif/src/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/runtime/middleware"
)

func Init() *gin.Engine {

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
		defer cancel()
		c.Set("ctx", ctx)
		c.Next()
	})

	uploader := &utils.GCSUploader{
		CDNUrl:          config.Config.Upload.CDN,
		BucketName:      config.Config.Upload.Bucket,
		CredentialsFile: config.Config.Upload.Credentials,
	}

	router.Use(func(c *gin.Context) {
		c.Set("uploader", uploader)
		c.Next()
	})

	//Cors
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if config.Config.Debug {
				return true
			}
			for _, o := range config.Config.Cors.Origins {
				if o == origin {
					return true
				}
			}
			return false
		},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	views.Init(router)

	//docs
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	router.GET("/docs", gin.WrapH(middleware.SwaggerUI(opts, nil)))
	router.GET("/swagger.yaml", gin.WrapH(http.FileServer(http.Dir("./docs"))))

	return router
}

func Serve() {
	router := Init()
	router.Run(fmt.Sprintf("0.0.0.0:%d", config.Config.Port))
}
