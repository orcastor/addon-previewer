package main

import (
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egin"

	"github.com/orcastor/orcas/core"
	"github.com/orcastor/orcas/rpc/middleware"
)

// ORCAS_CACHE=/opt/orcas_cache EGO_DEBUG=true EGO_LOG_EXTRA_KEYS=uid ORCAS_BASE=/opt/orcas ORCAS_DATA=/opt/orcas_disk ORCAS_SECRET=xxxxxxxx egoctl run --runargs --config=config.toml
// ORCAS_DOCKER_EXEC="docker exec -i <container_id>" ORCAS_CACHE=/opt/orcas_cache EGO_DEBUG=true EGO_LOG_EXTRA_KEYS=uid ORCAS_BASE=/opt/orcas ORCAS_DATA=/opt/orcas_disk ORCAS_SECRET=xxxxxxxx egoctl run --runargs --config=config.toml
// go run server.go --config=config.toml
func main() {
	core.InitDB()
	if err := ego.New().Serve(func() *egin.Component {
		server := egin.Load("server.http").Build()

		server.Use(middleware.Metrics())
		server.Use(middleware.CORS())
		server.Use(middleware.JWT())

		api := server.Group("/prvw/api")
		api.GET("/get", get)
		api.GET("/thumb/:to", thumb)
		return server
	}()).Run(); err != nil {
		elog.Panic("startup", elog.Any("err", err))
	}
}
