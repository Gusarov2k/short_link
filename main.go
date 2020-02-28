package main

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
)

type Link struct {
	OriginalLink string `json:"original_link" binding:"required"`
	ShortLink    string `json:"short_link" binding:"required"`
}

type LinkParams struct {
	OriginalLink string `form:"original_link" json:"original_link" validate:"url"`
}

var validate *validator.Validate

func main() {
	validate = validator.New()

	r := gin.Default()
	r.Use(location.Default())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/link", func(c *gin.Context) {
			var json LinkParams
			var link Link

			if err := c.ShouldBindJSON(&json); err != nil {
				c.JSON(422, gin.H{"error": err.Error()})
				return
			}
			if err := validate.Struct(&json); err != nil {
				c.JSON(422, gin.H{"error": err.Error()})
				return
			}

			b := make([]byte, 16)
			_, err := rand.Read(b)
			if err != nil {
				log.Fatal(err)
			}
			uuid := fmt.Sprintf("%s/short_link/%x-%x-%x-%x-%x",
				location.Get(c), b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
			link.ShortLink = uuid
			link.OriginalLink = json.OriginalLink
			c.JSON(200, gin.H{"data": link})
		})
		r.Run(":3002")
	}
	// https://stackoverflow.com/questions/46567672/what-is-the-best-way-to-handle-dynamic-subdomains-in-golang-with-gin-router
}
