package middleware

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		log.Printf("[INFO] %s %s", request.Method, request.RequestURI)

		next.ServeHTTP(writer, request)
	})
}

func HeaderNormalizerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/json; charset=utf-8")
		next.ServeHTTP(writer, request)
	})
}
