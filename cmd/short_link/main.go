package main

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"time"
)

// model
type UserLink struct {
	gorm.Model
	UserLink  string    `json:"user_link"`
	ShortLink string    `json:"short_link"`
	UserLog   []UserLog `gorm:"foreignkey:UserLinkID"`
}

type UserLog struct {
	gorm.Model
	Statistic  []Statistic `gorm:"foreignkey:UserStatisticID"`
	Ip         string      `json:"user_ip"`
	UserLinkID uint
}

type Statistic struct {
	Log             time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	UserStatisticID uint
}

type Link struct {
	UserLink  string `json:"original_link" binding:"required"`
	ShortLink string `json:"short_link" binding:"required"`
}

type LinkParams struct {
	OriginalLink string `form:"original_link" json:"original_link" validate:"url"`
}

var validate *validator.Validate

func main() {

	// db init
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=ivan dbname=short_link_development password=1234 sslmode=disable")

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	fmt.Printf("%s\n", err)
	db.LogMode(true)

	db.AutoMigrate(&UserLink{}, &UserLog{}, &Statistic{})

	//

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
			uuid := fmt.Sprintf("short=%x-%x-%x-%x-%x",
				b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

			var userLink = UserLink{UserLink: json.OriginalLink, ShortLink: uuid, UserLog: []UserLog{{Ip: c.ClientIP()}}}
			db.Save(&userLink)
			copier.Copy(&link, &userLink)
			c.JSON(200, gin.H{"data": link})
		})

		v1.GET("/short_link/:short", func(c *gin.Context) {
			name := c.Param("short")
			var user UserLink

			db.Where("short_link = ?", name).First(&user)

			if name == user.ShortLink {
				c.Redirect(301, user.UserLink)
			} else {
				c.JSON(404, gin.H{"error": "not find"})
			}
		})
		r.Run(":3002")
	}
}
