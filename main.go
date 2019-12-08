package main

import (
	"fmt"
	"net/http"

	"bloggist/api"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()
	fmt.Println("router init")
	be_api := router.Group("/api")
	{
		be_api.GET("/user/:name/blog", func(c *gin.Context) {
			fmt.Println("enter blogs")
			name := c.Param("name")
			c.JSON(200, gin.H{
				"blogs": []string{"JavaScript 学习笔记", "Vue 组件生命周期", "区块链作业好难顶"},
				"test":  3,
				"name":  name,
			})
		})
		be_api.POST("/user/:name/publish", api.PublishBlog)
		be_api.GET("/user/:name/blogs", api.GetBlogs)
		be_api.POST("/register", api.Register)
		be_api.POST("/login", api.Login)
		be_api.GET("/user/:name/blog/:blogid/like", api.LikeBlog)
		be_api.GET("/user/:name/blog/:blogid/delete", api.DeleteBlog)
		be_api.GET("/user/:name/blog/:blogid", api.GetBlog)
	}

	router.NoRoute(func(c *gin.Context) {
		fmt.Println("enter noroute")
		c.JSON(http.StatusNotFound, gin.H{
			"status": 404,
			"error":  "你来到了没有知识的荒野",
		})
	})

	router.Run(":8001")
}

func main() {
	Init()
}
