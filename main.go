package main

import (
	"fmt"
	"os"
	"text/template"
)

func usage() string {
	return fmt.Sprintf("usage: %s <template>", os.Args[0])
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, usage())
		os.Exit(1)
	}
	in := os.Args[1]

	t := template.New("template")

	t, err := t.Parse(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: parse template:", err)
		os.Exit(1)
	}

	if err := t.Execute(os.Stdout, nil); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: execute template:", err)
		os.Exit(1)
	}
}
