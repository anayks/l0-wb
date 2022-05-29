package app

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/nats-io/stan.go"
)

const (
	allowedHeaders = "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"
)

func (a *App) HandlerWebApp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)

	main := filepath.Join("html", "index.html")
	tmpl, err := template.ParseFiles(main)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func (a *App) HandlerGetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if r.Method != "POST" {
		fmt.Println("request method is not POST!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	defer r.Body.Close()

	result, err := a.Rep.Orders.GetJSONByID(data)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(result)
}

func (a *App) HandlerAddData(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error on createing data"))
		}
	}()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)

	if r.Method != "POST" {
		fmt.Println("request method is not POST!")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Printf("error on reading body: %v", err)
		return
	}
	r.Body.Close()

	sc := *a.Sc

	sc.Publish("json-parser", data)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("element created"))
}

func (a *App) HandleJSON(m *stan.Msg) {
	ctx := context.Background()
	if err := a.Rep.Orders.CreateByJSON(ctx, m.Data); err != nil {
		fmt.Printf("error on creating order: %v", err)
	}
}
