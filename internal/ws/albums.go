package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milantracy/playground/api/ws"
)

var albums = []ws.Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func GetAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func AddAlbum(c *gin.Context) {
	var newAlbum ws.Album

	if err := c.BindJSON(&newAlbum); err == nil {
		albums = append(albums, newAlbum)
		c.IndentedJSON(http.StatusCreated, newAlbum)
	}
}
