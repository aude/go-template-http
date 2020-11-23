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

GET /environ to see available env vars

examples:

	$ curl http://localhost/
	$ curl http://localhost/environ
	$ curl -d 'The value of $PATH is: {{getenv "PATH"}}' http://localhost/
`)
}

var funcMap = template.FuncMap{
	"getenv": os.Getenv,
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

		tmpl, err := template.New("random string").Funcs(funcMap).Parse(string(body))
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
		fmt.Fprintln(w)
	})

	http.HandleFunc("/environ", func(w http.ResponseWriter, r *http.Request) {
		// if not GET, reply with usage instructions
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintln(w, usage())
			return
		}

		for _, env := range os.Environ() {
			fmt.Fprintln(w, env)
		}
	})

	log.Fatal(http.ListenAndServe(":80", nil))
}
