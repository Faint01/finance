package main

import (
	docs "github.com/Faint01/finance/docs"
	hendler "github.com/Faint01/finance/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@BasePath	/api/v1
//	@Title			Test Swagger API
//	@Description	This is a sample server.
//	@Version		1.0

func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/api/v1/all", hendler.GetAll)

	r.GET("/api/v1/finance/:id", hendler.IdSearch)

	r.PUT("/api/v1/updatefin/:id", hendler.Updatefin)

	r.POST("/api/v1/addfin", hendler.Postfinc)

	r.DELETE("/api/v1/removefin/:id", hendler.RemoveRecord)

	r.Run(":8080")
}
