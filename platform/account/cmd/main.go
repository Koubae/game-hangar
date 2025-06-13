package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))

	var accounts gin.Accounts
	//accounts = gin.Accounts{
	//	"foo":  "bar",
	//	"manu": "123",
	//}
	accounts = gin.Accounts{
		"admin": "admin",
		"foo":   "bar",
		"manu":  "123",
	}
	//authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	//	"foo":  "bar", // user:foo password:bar
	//	"manu": "123", // user:manu password:123
	//}))
	authorized := r.Group("/", gin.BasicAuth(accounts))

	/* example curl for /admin with basicauth header
		Zm9vOmJhcg== is base64("foo:bar")

		YWRtaW46YWRtaW4= base64("admin:admin")

	curl -X POST \
	http://localhost:8080/admin \
	-H 'authorization: Basic YWRtaW46YWRtaW4=' \
	-H 'content-type: application/json' \
	-d '{"value":"admin"}'

		curl -X POST \
		http://localhost:8080/admin \
		-H 'authorization: Basic Zm9vOmJhcg==' \
		-H 'content-type: application/json' \
		-d '{"value":"bar"}'

		curl -X POST \
		http://localhost:8080/admin \
		-H 'authorization: Basic bWFudToxMjM=' \
		-H 'content-type: application/json' \
		-d '{"value":"123"}'

		curl -X POST \
		http://localhost:8080/admin \
		-H 'authorization: Basic blabla' \
		-H 'content-type: application/json' \
		-d '{"value":"123"}'

	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	/*
		curl -X POST 'http://localhost:8080/admin/user_1?password=pass' -H 'authorization: Basic YWRtaW46YWRtaW4='
	*/
	authorized.POST("admin/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		password, exists := c.GetQuery("password")
		if !exists {
			c.String(400, "Password parameter is required")
			return
		}

		accounts[name] = password

		credential := fmt.Sprintf("%s:%s", name, password)
		encoded := base64.StdEncoding.EncodeToString([]byte(credential))

		requestSample := fmt.Sprintf(`curl -X POST \
		http://localhost:8080/admin \
		-H 'authorization: Basic %s' \
		-H 'content-type: application/json' \
		-d '{"value":"%s"}'`, encoded, password)

		//c.JSON(http.StatusOK, gin.H{"status": "ok", "credential": encoded, "accounts": accounts, "requestSample": requestSample})
		fmt.Println(gin.DefaultWriter, accounts)
		c.String(http.StatusOK, requestSample)
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
