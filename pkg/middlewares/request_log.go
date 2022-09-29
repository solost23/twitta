package middlewares

import (
	"Twitta/pkg/utils"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
	"time"
)

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		RequestInLog(c)
		c.Next()
		RequestOutLog(c)
	}
}

func RequestInLog(c *gin.Context) {
	traceId := utils.UUID()
	c.Set("startTime", time.Now())
	c.Set("traceId", traceId)

	uri := c.Request.RequestURI
	ua := c.Request.Header.Get("User-Agent")
	if !strings.Contains(uri, "login") {
		logContent := map[string]interface{}{
			"ua":      ua,
			"uri":     uri,
			"method":  c.Request.Method,
			"args":    c.Request.PostForm,
			"from":    c.ClientIP(),
			"traceId": traceId,
		}
		if c.ContentType() == "application/json" {
			bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			logContent["body"] = string(bodyBytes)
		}
		go zap.S().Info(logContent)
	}
}

func RequestOutLog(c *gin.Context) {
	// after request
	endTime := time.Now()
	st, _ := c.Get("startTime")

	isJSON := false
	contentTypes := c.Writer.Header()["Content-Type"]
	for _, ct := range contentTypes {
		if strings.Index(ct, "application/json") != -1 {
			isJSON = true
			break
		}
	}

	startExecTime, _ := st.(time.Time)
	traceId := c.MustGet("traceId")
	userId := c.Value("userId")
	logContent := map[string]interface{}{
		"userId":    userId,
		"uri":       c.Request.RequestURI,
		"method":    c.Request.Method,
		"args":      c.Request.PostForm,
		"from":      c.ClientIP(),
		"traceId":   traceId,
		"proc_time": endTime.Sub(startExecTime).Seconds(),
	}
	if isJSON {
		response, _ := c.Get("response")
		logContent["response"] = fmt.Sprintf("%#v", response)
	}
	go zap.S().Info(logContent)
}
