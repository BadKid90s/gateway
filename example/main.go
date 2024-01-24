package main

import (
	"knife"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/a", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello a "))
	})

	chain := knife.NewChain(mux)

	chain.Use(func(context *knife.Context) {
		start := time.Now()
		defer func() {
			duration := time.Now().Sub(start)
			log.Printf("process all middleware,time consumption %s ", duration)
		}()
		context.Next()
	})

	chain.Use(func(context *knife.Context) {
		log.Printf("logger middleware ,path:%s ", context.Req.URL.Path)
		context.Writer.Header().Set("token", "123")

		context.Next()
	})

	chain.UseMatcher(func(response knife.HttpResponse, request knife.HttpRequest) bool {
		header := response.Header()
		return header.Get("token") != ""
	}, func(context *knife.Context) {
		log.Printf("token middleware,token:%s ", context.Writer.Header().Get("token"))
		context.Next()
	})

	_ = http.ListenAndServe(":8080", chain)
}
