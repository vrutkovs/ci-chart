package ui

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/vrutkovs/ci-chart/pkg/event"
	"k8s.io/klog/v2"
)

func Run(store event.Store, port uint16) {
	r := mux.NewRouter()
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		path = "."
	}
	r.Handle("/", http.FileServer(http.Dir(fmt.Sprintf("%s/static", path))))
	r.HandleFunc("/data.json", store.JSONHandler)
	klog.Infof(fmt.Sprintf("Listening on :%d", port))
	klog.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
