package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type Request map[string]interface{}
type Response interface{}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {

}

func NotFound(w http.ResponseWriter, r *http.Request) {

}

func Exec(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	vm := otto.New()

	concierge := map[string]string{
		"JSON":  "application/json; charset=utf-8",
		"PLAIN": "text/plain; charset=utf-8",
	}

	request := make(map[string]interface{}, 2)
	request["URL"] = r.URL.String()
	request["params"] = ps.ByName

	next := func(statusCode int, contentType string, body interface{}) interface{} {
		w.WriteHeader(statusCode)
		w.Header().Set("content-type", contentType)
		return body
	}

	vm.Set("concierge", concierge)

	vm.Set("request", request)

	vm.Set("next", next)

	script, err := vm.Compile("./server.js", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	value, err := vm.Run(script)
	if err != nil {
		fmt.Println(err.Error())
	}

	if value.IsString() {
		w.Write([]byte(value.String()))
	} else {
		w.WriteHeader(500)
		w.Header().Set("content-type", "application/json; charset=utf-8")
		w.Write([]byte(`{"statusCode":500, "code": 500.1, "error": "Invalid response"}`))
	}

}

func main() {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(NotFound)
	router.MethodNotAllowed = http.HandlerFunc(MethodNotAllowed)

	router.GET("/execute/:name", Exec)

	handler := fasthttpadaptor.NewFastHTTPHandler(router)

	fasthttp.ListenAndServe(":8000", handler)

}
