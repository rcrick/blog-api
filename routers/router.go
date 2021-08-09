package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrick/blog-api/middleware/jwt"
	"github.com/rcrick/blog-api/pkg/setting"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "github.com/rcrick/blog-api/docs"
	v1 "github.com/rcrick/blog-api/routers/api/v1"
)

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @termsOfService http://swagger.io/terms/

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/auth", v1.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	api.Use(jwt.JWT())
	{
		api.GET("/tags", v1.GetTags)
		api.POST("/tags", v1.AddTag)
		api.PUT("/tag/:id", v1.EditTag)
		api.DELETE("/tag/:id", v1.DeleteTag)

		api.GET("/articles", v1.GetArticles)
		api.GET("/article/:id", v1.GetArticle)
		api.POST("/article", v1.AddArticle)
		api.PUT("/article/:id", v1.EditArticle)
		api.DELETE("/article/:id", v1.DeleteArticle)
	}

	return r
}
