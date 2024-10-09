package middleware

import (
	"log"
	"mime"
	"net/http"
)

// Logging 打印请求参数
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("recv a %s request from %s", request.Method, request.RemoteAddr)
		next.ServeHTTP(writer, request)
	})
}

// Validating 校验请求头的 Content-Type 是否为 application/json
func Validating(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		contentType := request.Header.Get("Content-Type")
		mediatype, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if mediatype != "application/json" {
			http.Error(writer, "invalid Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
