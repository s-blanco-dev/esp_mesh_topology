package main

import (
	"mesh_topology/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/static", "./static")
	r.POST("/v1/update", handler.PostMeshData)

	r.Run(":8080")
}
