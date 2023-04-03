package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jeffleon1/transaction-ms/pkg/health"
	"github.com/jeffleon1/transaction-ms/pkg/swagger"
	"github.com/jeffleon1/transaction-ms/pkg/transactions/infrastructure"
)

type Router interface {
	Run(addr ...string) error
}

func NewRouter(routes RoutesGroup) Router {
	route := gin.Default()
	public := route.Group("/api/stori/v1/public")
	routes.Transaction.PublicRoutes(public)
	routes.Health.RegisterRoutes(public)
	routes.Swagger.RegisterRoutes(public)

	return route
}

type RoutesGroup struct {
	Transaction *infrastructure.TransactionRoutes
	Health      *health.Routes
	Swagger     *swagger.Routes
}
