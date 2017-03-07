package main

import (
	"html"
	"strings"
)

const errormessage = `<div class="alert alert-danger" role="alert">
  <span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span>
  <span class="sr-only">Error:</span>
  $MESSAGE$
</div>`

func BuildMessage(template string, message string) string {
	message = html.EscapeString(message)
	return strings.Replace(template, "$MESSAGE$", message, -1)
}
