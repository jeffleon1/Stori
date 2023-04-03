package swagger

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffleon1/transaction-ms/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
