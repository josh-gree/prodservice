package main

import (
	"fmt"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/op/go-logging"
	"net/http"
	"bytes"
	"os"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{pid} %{longfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
)

type Job struct{
	Data []float64 `json:"data"`
}

type Result struct{
	Out float64 `json:"out"`
	Service string `json:"service"`
}

func main(){

	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend1, backend1Formatter)

	Listen()
}

func Send(r Result) error {
	log.Info("Sending result: prod")
	data, err := json.Marshal(r)

	// How to deal with name issues?? service only needs to know public location
	publichost := "localhost:7000"
	_, err = http.Post(fmt.Sprintf("http://%s/result/",publichost),"application/json",bytes.NewBuffer(data))
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}


func Recivejob(c echo.Context) error {
	log.Info("Recieved Job: prod")
	j := Job{}
	err := c.Bind(&j)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("Data : ",j.Data)

	result := Prod(j)

	Send(result)

	return nil
}

func Prod(j Job) Result {
	log.Info("Doing computation: prod")
	sum := 1.0
	for _, v := range j.Data {
		sum *= v
	}
	res := Result{Out:sum,Service:"sum"}
	return res
}

func Listen() {
	log.Info("Starting to Listen prod")
	e := echo.New()

	e.POST("/", Recivejob)
	e.Start(":9000")
}
