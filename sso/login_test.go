package _sso

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type User struct {
	Name string
	Age  int
}

func TestBizServer(t *testing.T) {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	r.GET("/login", func(c *gin.Context) {
		var user User
		err := c.BindJSON(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.String(200, "login")
	})
	r.Use(Authorization())
	r.GET("/profile", Profile)
	r.Run(":8080")
}

func Authorization() func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO: 鉴权逻辑
		c.JSON(http.StatusOK, gin.H{"message": "Authorization"})
		c.Next()
	}
}

func Profile(c *gin.Context) {
	c.JSON(200, User{Name: "test", Age: 18})
}
