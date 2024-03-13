package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LogLayout Log layout
type LogLayout struct {
	Time      time.Time
	Metadata  map[string]interface{} // Stores custom raw data
	Path      string                 // Access path
	Query     string                 // Carries query
	Body      string                 // Carries the body data
	IP        string                 // ip address
	UserAgent string                 // Agent
	Error     string                 // Error
	Cost      time.Duration          // Cost time
	Source    string                 // Source
}

type Logger struct {
	Filter func(c *gin.Context) bool

	FilterKeyword func(layout *LogLayout) bool

	AuthProcess func(c *gin.Context, layout *LogLayout)

	Print func(LogLayout)

	Source string
}

func (l Logger) SetLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		var body []byte
		if l.Filter != nil && !l.Filter(c) {
			body, _ = c.GetRawData()

			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		c.Next()
		cost := time.Since(start)
		layout := LogLayout{
			Time:      time.Now(),
			Path:      path,
			Query:     query,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Error:     strings.TrimRight(c.Errors.ByType(gin.ErrorTypePrivate).String(), "\n"),
			Cost:      cost,
			Source:    l.Source,
		}
		if l.Filter != nil && !l.Filter(c) {
			layout.Body = string(body)
		}

		l.AuthProcess(c, &layout)
		if l.FilterKeyword != nil {

			l.FilterKeyword(&layout)
		}

		l.Print(layout)
	}
}

func DefaultLogger() gin.HandlerFunc {
	return Logger{
		Print: func(layout LogLayout) {
			v, _ := json.Marshal(layout)
			fmt.Println(string(v))
		},
		Source: "ZC",
	}.SetLoggerMiddleware()
}
