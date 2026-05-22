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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")

		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			hasValidToken := (err == nil && api.IsValidToken(cookie.Value))

			if (r.URL.Path == "/" || r.URL.Path == "/index.html") && !hasValidToken {
				http.Redirect(w, r, "login.html", http.StatusFound)
				return
			}

			if hasValidToken && r.URL.Path == "/login.html" {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
		}

		fs.ServeHTTP(w, r)
	})

	port := "7540"
	if env := os.Getenv("TODO_PORT"); env != "" {
		port = env
	}
	addr := ":" + port

	fmt.Printf("starting server on the port %s\n", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Error %s", err.Error())
	}

}
