package main

import (
	"crypto/rand"
	"fmt"
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
	gin.ForceConsoleColor()
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/link", func(c *gin.Context) {
			var json LinkParams
			var link Link
			log.Println(c.Request.URL.Path)

			if err := c.ShouldBindJSON(&json); err != nil {
				c.JSON(422, gin.H{"error": err.Error()})
				return
			}
			if err := validate.Struct(&json); err != nil {
				c.JSON(422, gin.H{"error": err.Error()})
				return
			}
			link.OriginalLink = json.OriginalLink
			b := make([]byte, 16)
			_, err := rand.Read(b)
			if err != nil {
				log.Fatal(err)
			}
			uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
				b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
			link.ShortLink = uuid

			c.JSON(200, gin.H{"data": link})
		})
		r.Run(":3002")
	}

}
