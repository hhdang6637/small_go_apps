package htmlHelper

import (
	"fmt"
	"net/http"
)

const (
	docBegin  = "<!DOCTYPE html>"
	htmlBegin = "<html>"
	htmlEnd   = "</html>"
	headBegin = "<head>"
	headEnd   = "</head>"

	cssStyle = `
<style>
table {
font-family: arial, sans-serif;
border-collapse: collapse;
width: 100%;
}

td, th {
border: 1px solid #dddddd;
text-align: left;
padding: 8px;
}

tr:nth-child(even) {
background-color: #dddddd;
}
</style>
	`

	bodyBegin = "<body>"
	bodyEnd   = "</body>"
)

// HTMLDocOpen write standar html text
func HTMLDocOpen(w http.ResponseWriter) {
	fmt.Fprintln(w, docBegin)
	fmt.Fprintln(w, htmlBegin)
	fmt.Fprintln(w, headBegin)
	fmt.Fprintln(w, cssStyle)
	fmt.Fprintln(w, headEnd)
	fmt.Fprintln(w, bodyBegin)
}

// HTMLDocClose write ending body and html tags
func HTMLDocClose(w http.ResponseWriter) {
	fmt.Fprintf(w, bodyEnd)
	fmt.Fprintf(w, htmlEnd)
}

// TableDocOpen write standar table tags
func TableDocOpen(w http.ResponseWriter, tHeaders []string) {
	fmt.Fprintln(w, "<table><tr>")
	for _, header := range tHeaders {
		fmt.Fprintf(w, "<th>%s</th>\n", header)
	}
	fmt.Fprintln(w, "</tr>")

}

// TableDocClose write table ending tag
func TableDocClose(w http.ResponseWriter) {
	fmt.Fprintln(w, "</table>")
}
