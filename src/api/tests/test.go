package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	router := gin.Default()
	// 封装request和response到gin.Context中
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// gin路由来自httprouter库
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")
		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	// body的传输格式包括application/json, application/x-www-form-urlencoded, application/xml, multipart/formdata
	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(http.StatusOK, gin.H{
			"status": gin.H{
				"status_code": http.StatusOK,
				"status":      http.StatusText(http.StatusOK),
			},
			"message": message,
			"nick":    nick,
		})
	})
	router.PUT("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "1")
		name := c.PostForm("name")
		message := c.PostForm("message")
		fmt.Printf("id : %s, page : %s, name : %s, message : %s\n", id, page, name, message)
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusText(http.StatusOK),
		})
	})
	// 上传文件
	router.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		fmt.Println(name)
		file, header, err := c.Request.FormFile("upload")
		if err != nil {
			c.String(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		filename := header.Filename
		fmt.Println(file, err, filename)
		out, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		// 保存上传的文件到指定位置
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		c.String(http.StatusOK, http.StatusText(http.StatusOK))
	})
	// 上传多个文件
	router.POST("/multi/upload", func(c *gin.Context) {
		err := c.Request.ParseMultipartForm(20000)
		if err != nil {
			log.Fatal(err)
		}
		formData := c.Request.MultipartForm

		files := formData.File["upload"]
		for i := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			out, err := os.Create(files[i].Filename)
			defer out.Close()
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(out, file)
			if err != nil {
				log.Fatal(err)
			}
		}
		c.String(http.StatusOK, http.StatusText(http.StatusOK))
	})
	// 判断content-type
	router.POST("/login", func(c *gin.Context) {
		var user User
		var err error
		contentType := c.Request.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			err = c.BindJSON(&user)
		case "application/x-www-form-urlencoded":
			err = c.BindWith(&user, binding.Form)
		}
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"user":   user.Username,
			"passwd": user.Passwd,
			"age":    user.Age,
		})
	})
	// 自动判断content-type
	router.POST("/logins", func(c *gin.Context) {
		var user User
		err := c.Bind(&user)
		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"user":   user.Username,
			"passwd": user.Passwd,
			"age":    user.Age,
		})
	})
	// 多格式渲染
	router.GET("render", func(c *gin.Context) {
		contentType := c.DefaultQuery("content_type", "json")
		if contentType == "json" {
			c.JSON(http.StatusOK, gin.H{
				"status": http.StatusText(http.StatusOK),
			})
		} else {
			c.XML(http.StatusOK, gin.H{
				"status": http.StatusText(http.StatusOK),
			})
		}
	})
	// 重定向路由
	router.GET("/redirect/google", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "http://www.google.com")
	})
	// 分组路由
	v1 := router.Group("v1")
	v1.GET("/login", func(context *gin.Context) {
		context.String(http.StatusOK, http.StatusText(http.StatusOK))
	})
	v2 := router.Group("v2")
	v2.GET("/login", func(context *gin.Context) {
		context.String(http.StatusOK, http.StatusText(http.StatusOK))
	})
	// 全局中间件
	router.Use(MiddleWare())
	{
		router.GET("/middle", func(context *gin.Context) {
			request := context.MustGet("client_request").(string)
			req, _ := context.Get("request")
			context.JSON(http.StatusOK, gin.H{
				"middle_request": request,
				"request":        req,
			})
		})
	}
	// 单个路由中间件
	router.GET("/sigle_middle", MiddleWare(), func(context *gin.Context) {
		request := context.MustGet("client_request").(string)
		req, _ := context.Get("request")
		context.JSON(http.StatusOK, gin.H{
			"middle_request": request,
			"request":        req,
		})
	})
	// 群组路由中间件
	authorized := router.Group("/")
	authorized.Use(MiddleWare())
	{
		router.GET("/group_middle", func(context *gin.Context) {
			request := context.MustGet("client_request").(string)
			req, _ := context.Get("request")
			context.JSON(http.StatusOK, gin.H{
				"middle_request": request,
				"request":        req,
			})
		})
	}
	// 中间价可以用来记录log,追踪请求,错误处理,接口鉴权
	router.GET("/auth/signin", func(c *gin.Context) {
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    "123",
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)
		c.String(http.StatusOK, "Login successful")
	})

	router.GET("/home", AuthMiddleWare(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "home"})
	})
	// 异步协程
	router.GET("/sync", func(context *gin.Context) {
		time.Sleep(5 * time.Second)
		context.String(http.StatusOK, "success")
	})
	router.GET("/async", func(context *gin.Context) {
		contextCopy := context.Copy()
		go func() {
			fmt.Println("Start sleep")
			time.Sleep(5 * time.Second)
			fmt.Println("sleep end")
			fmt.Println("Done, in url" + contextCopy.Request.URL.Path)
		}()
		context.String(http.StatusOK, "success")
	})

	// 自定义router
	// 与net/http结合
	/**
	router := gin.Default()
	http.ListenAndServe(":8080", router)
	*/
	/**
	router := gin.Default()
	server := &http.Server{
		Addr:":8080",
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20
	}
	server.ListenAndServe()
	*/
	router.Run() // listen and serve on 0.0.0.0:8080
}

type User struct {
	Username string `form:"username" json:"username" binding:"required"`
	Passwd   string `form:"passwd" json:"passwd" binding:"required"`
	Age      int    `form:"age" json:"age"`
}

func MiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("before middleware")
		context.Set("request", "client_request")
		context.Next()
		fmt.Println("after middleware")
	}
}
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie("session_id"); err == nil {
			value := cookie.Value
			fmt.Println(value)
			if value == "123" {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}
}
