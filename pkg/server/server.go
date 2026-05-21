package server

import (
	"fmt"
	"go_final_project/pkg/api"
	"net/http"
	"os"
)

const webDir = "./web"

func Run() {
	api.Init()

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	port := "7540"
	if env := os.Getenv("TODO_PORT"); env != "" {
		port = env
	}
	addr := ":" + port

	fmt.Printf("starting server on the port %s\n", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Ошибка %s", err.Error())
	}

}
