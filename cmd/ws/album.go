package main

import (
	"github.com/gin-gonic/gin"
	"github.com/milantracy/playground/internal/ws"
)

func main() {
	router := gin.Default()
	router.GET("/albums", ws.GetAlbums)
	router.GET("/albums/:id", ws.GetAlbumByID)
	router.POST("/albums", ws.AddAlbum)

	router.Run("localhost:8080")
}
