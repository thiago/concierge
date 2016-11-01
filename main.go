package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func Execute(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	errch := make(chan error, 1)
	cmd := exec.Command("node", "./wrapper.js")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		errch <- cmd.Wait()
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			w.Write([]byte(line))
		}
	}()

	select {
	case <-time.After(time.Millisecond * 80):
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Println("Timeout hit..")
		return
	case err := <-errch:
		if err != nil {
			log.Println("traceroute failed:", err)
		}
	}
}

func main() {

	port := ":8000"

	log.Println("Server started on port:", port)

	router := httprouter.New()

	router.GET("/:name", Execute)

	handler := fasthttpadaptor.NewFastHTTPHandler(router)

	if err := fasthttp.ListenAndServe(port, handler); err != nil {
		log.Println(err.Error())
	}

}
