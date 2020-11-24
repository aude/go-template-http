package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

func usage() string {
	return fmt.Sprintf(`POST Go templates to me, I'll render them :D

you can use my environment variables as well, through the "getenv" function

list available env vars with the "environ" function

examples:

	$ curl http://localhost:8080/
	$ curl -d 'The value of $PATH is: {{getenv "PATH"}}' http://localhost:8080/
	$ curl -d '{{environ}}' http://localhost:8080/
	$ curl -d '{{range environ}}{{println .}}{{end}}' http://localhost:8080/
`)
}

var funcMap = template.FuncMap{
	"environ": os.Environ,
	"getenv":  os.Getenv,
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var err error

		// if not POST, reply with usage instructions
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w, usage())
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "could not read request body:", err)
			return
		}

		// if empty request body, reply with usage instructions
		if string(body) == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, usage())
			return
		}

		tmpl, err := template.New("request").Funcs(funcMap).Parse(string(body))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "could not parse template:", err)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "could not execute template:", err)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
