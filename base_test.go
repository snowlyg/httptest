package httptest

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type Request struct {
	Message string `json:"message" form:"message" uri:"message"`
}

// GinHandler Create add /example route to gin engine
func GinHandler(r *gin.Engine) *gin.Engine {

	// Add route to the gin engine
	r.GET("/example", func(c *gin.Context) {
		var req Request
		if errs := c.ShouldBind(&req); errs != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		message := "pong"
		if req.Message != "" {
			message = req.Message
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "OK",
			"data": gin.H{
				"message": message,
			},
		})
	})

	// Add route to the gin engine
	r.GET("/mutil", func(c *gin.Context) {
		var req Request
		if errs := c.ShouldBind(&req); errs != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		message := "pong"
		if req.Message != "" {
			message = req.Message
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "OK",
			"data": []gin.H{
				{"message": message},
				{"message": message},
			},
		})
	})

	// Add route to the gin engine
	r.POST("/example", func(c *gin.Context) {
		var req Request
		if errs := c.ShouldBindJSON(&req); errs != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "FAIL",
			})
			return
		}
		message := "pong"
		if req.Message != "" {
			message = req.Message
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "OK",
			"data": gin.H{
				"message": message,
			},
		})
	})

	// Add route to the gin engine
	r.POST("/upload", func(c *gin.Context) {
		_, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "FAIL",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "操作成功",
		})
	})

	type RequestId struct {
		Id uint `json:"id" form:"id" uri:"id"`
	}

	// Add route to the gin engine
	r.DELETE("/example/:id", func(c *gin.Context) {
		var req RequestId
		if errs := c.ShouldBindUri(&req); errs != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "OK",
			"data": gin.H{
				"id": req.Id,
			},
		})
	})

	// Add route to the gin engine
	r.POST("login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "OK",
			"data": gin.H{
				"AccessToken": "EIIDFJDIKFJJIdfdkfk.uisdifsdfisdouf",
				"user": gin.H{
					"id": 1,
				},
			},
		})
	})

	// Add route to the gin engine
	r.GET("logout", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "操作成功",
		})
	})

	// Add route to the gin engine
	r.GET("header", func(c *gin.Context) {
		c.GetHeader("Authorization")
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "操作成功",
			"data": gin.H{
				"Authorization": c.GetHeader("Authorization"),
			},
		})
	})

	// return gin engine with newly added route
	return r
}

func TestInstance(t *testing.T) {
	engine := gin.New()
	// Create httpexpect instance
	client := Instance(t, GinHandler(engine))
	client.GET("/example", NewResponses(http.StatusOK, "OK", Responses{{Key: "message", Value: "pong"}}))
	client.DELETE("/example/1", NewResponses(http.StatusOK, "OK", Responses{{Key: "id", Value: 1}}))
}

func TestNewWithQueryObjectParamFunc(t *testing.T) {
	engine := gin.New()
	// Create httpexpect instance
	client := Instance(t, GinHandler(engine))
	pageKeys := Responses{{Key: "message", Value: "message"}}
	client.GET("/example", NewResponses(http.StatusOK, "OK", pageKeys), NewWithQueryObjectParamFunc(map[string]interface{}{"message": "message"}))
}

func TestNewNewWithJsonParamFunc(t *testing.T) {
	engine := gin.New()
	// Create httpexpect instance
	client := Instance(t, GinHandler(engine))
	client.POST("/example", NewResponses(http.StatusOK, "OK", Responses{{Key: "message", Value: "message"}}), NewWithJsonParamFunc(map[string]interface{}{"message": "message"}))
	client.POST("/example", NewResponses(http.StatusOK, "OK", Responses{{Key: "message", Value: "pong"}}), NewWithJsonParamFunc(map[string]interface{}{"message": ""}))
}

func TestNewResponses(t *testing.T) {
	engine := gin.New()
	// Create httpexpect instance
	client := Instance(t, GinHandler(engine))

	client.GET("/example", NewResponses(http.StatusOK, "OK", Responses{{Key: "message", Value: "pong"}}))
	client.GET("/mutil", NewResponses(http.StatusOK, "OK", Responses{{Key: "message", Value: "pong"}}, Responses{{Key: "message", Value: "pong"}}))
	client.POST("/example", NewResponses(http.StatusBadRequest, "FAIL"))
}

func TestNewResponsesWithLength(t *testing.T) {
	engine := gin.New()
	// Create httpexpect instance
	client := Instance(t, GinHandler(engine))
	res := []Responses{{{Key: "message", Value: "pong"}}, {{Key: "message", Value: "pong"}}}
	client.GET("/mutil", NewResponsesWithLength(http.StatusOK, "OK", res, 2))
}

func TestNewWithFileParamFunc(t *testing.T) {
	engine := gin.New()
	// Create httpexpect instance
	client := Instance(t, GinHandler(engine))
	name := "test_img.jpg"
	fh, _ := os.Open("./" + name)
	defer fh.Close()
	client.UPLOAD("/upload", SuccessResponse, NewWithFileParamFunc([]File{{Key: "file", Path: name, Reader: fh}}))
}

func TestCheckStatus(t *testing.T) {
	status := checkStatus(NewResponses(http.StatusOK, "OK"))
	if status != http.StatusOK {
		t.Errorf("checkStatus want %d but get %d", http.StatusOK, status)
	}
	status = checkStatus(NewResponses(http.StatusBadRequest, "FAIL"))
	if status != http.StatusBadRequest {
		t.Errorf("checkStatus want %d but get %d", http.StatusBadRequest, status)
	}
	status = checkStatus(Responses{})
	if status != http.StatusOK {
		t.Errorf("checkStatus want %d but get %d", http.StatusOK, status)
	}
}

func TestLogin(t *testing.T) {
	engine := gin.New()
	// Create httpexpect instance
	client := Instance(t, GinHandler(engine))
	id, err := client.Login("/login", NewResponses(http.StatusOK, "OK", Responses{{Key: "AccessToken", Value: "EIIDFJDIKFJJIdfdkfk.uisdifsdfisdouf"}, {Key: "user", Value: Responses{{Key: "id", Value: 1}}}}))
	if err != nil {
		t.Error(err.Error())
	}
	if id == 0 {
		t.Error("id is 0")
	}
	client.GET("/header", NewResponses(http.StatusOK, "操作成功", Responses{{Key: "Authorization", Value: "Bearer EIIDFJDIKFJJIdfdkfk.uisdifsdfisdouf"}}))
}

func TestLogout(t *testing.T) {
	engine := gin.New()
	client := Instance(t, GinHandler(engine))
	client.Logout("/logout", SuccessResponse)
}
