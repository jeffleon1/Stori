package swagger

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

type Routes struct{}

func NewSwaggerDocsRoutes() *Routes {
	return &Routes{}
}

func (r *Routes) RegisterRoutes(group *gin.RouterGroup) {
	group.GET(
		"/docs/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.InstanceName(docs.SwaggerInfo.InstanceName()),
		),
	)
}
