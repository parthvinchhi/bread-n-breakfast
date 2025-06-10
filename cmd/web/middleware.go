package main

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"
)

// NoSurfGin adapts NoSurf for Gin usage
func NoSurfGin(inProduction bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create the NoSurf handler
		csrfHandler := nosurf.New(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Call the next Gin handler
			c.Request = r
			c.Writer = &ginResponseWriter{ResponseWriter: w, ginWriter: c.Writer}
			c.Next()
		}))
		csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   inProduction,
			SameSite: http.SameSiteLaxMode,
		})
		csrfHandler.ServeHTTP(c.Writer, c.Request)
	}
}

// ginResponseWriter is used to bridge between http.ResponseWriter and gin.ResponseWriter
type ginResponseWriter struct {
	ResponseWriter http.ResponseWriter
	ginWriter      gin.ResponseWriter
}

func (g *ginResponseWriter) Write(b []byte) (int, error) {
	return g.ginWriter.Write(b)
}

func SessionLoad(session *scs.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler := session.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// update Gin's request with the one returned from session
			c.Request = r
			c.Next()
		}))
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
