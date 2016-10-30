package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ddliu/motto"
	"github.com/julienschmidt/httprouter"
	"github.com/robertkrimen/otto"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type Request map[string]interface{}
type Response interface{}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(415)
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.Write([]byte(`{"statusCode":415, "code": 415.1, "error": "Resource not work with method"}`))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.Write([]byte(`{"statusCode":404, "code": 404.1, "error": "Resource not found"}`))
}

func VMJavaScript(w http.ResponseWriter, r *http.Request, ps httprouter.Params, chanDone chan<- error) {

	vm := motto.New()

	concierge := map[string]string{
		"JSON":  "application/json; charset=utf-8",
		"PLAIN": "text/plain; charset=utf-8",
	}

	request := make(map[string]interface{}, 2)
	request["URL"] = r.URL.String()
	request["params"] = ps.ByName

	response := func(statusCode int, contentType string, body interface{}) interface{} {
		w.WriteHeader(statusCode)
		w.Header().Set("content-type", contentType)
		return body
	}

	vm.Set("concierge", concierge)

	src := fmt.Sprintf("./%s.js", ps.ByName("file"))

	_, err := vm.Compile(src, nil)
	if err != nil {
		log.Println(err.Error())
	}

	exports, err := vm.Run(src)
	if err != nil {
		log.Println(err.Error())
	}

	var value otto.Value

	if exports.IsFunction() {

		value, err = exports.Call(exports, request, response)
		if err != nil {
			log.Println(err.Error())
		}

	} else {

		fn, err := exports.Object().Get(ps.ByName("fn"))
		if err != nil {
			log.Println(err.Error())
		}

		value, err = fn.Call(fn, request, response)
		if err != nil {
			log.Println(err.Error())
		}

	}

	fmt.Println(value)

	if value.IsString() {
		w.Write([]byte(value.String()))
	} else {
		w.WriteHeader(500)
		w.Header().Set("content-type", "application/json; charset=utf-8")
		w.Write([]byte(`{"statusCode":500, "code": 500.1, "error": "Invalid response"}`))
	}

	chanDone <- err

}

func Execute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	done := make(chan error)

	go VMJavaScript(w, r, ps, done)

	select {
	case <-time.After(3 * time.Second):
		w.WriteHeader(500)
		w.Header().Set("content-type", "application/json; charset=utf-8")
		w.Write([]byte(`{"statusCode":500, "code": 500.2, "error": "Timout in process"}`))

	case err := <-done:
		if err != nil {
			w.WriteHeader(500)
			w.Header().Set("content-type", "application/json; charset=utf-8")
			w.Write([]byte(`{"statusCode":500, "code": 500.3, "error": ` + err.Error() + `}`))
		}
	}

}

func main() {

	port := ":8000"

	log.Println("Server started on port:", port)

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(NotFound)
	router.MethodNotAllowed = http.HandlerFunc(MethodNotAllowed)

	router.GET("/:file/:fn/:name", Execute)

	handler := fasthttpadaptor.NewFastHTTPHandler(router)

	if err := fasthttp.ListenAndServe(port, handler); err != nil {
		log.Println(err.Error())
	}

}
