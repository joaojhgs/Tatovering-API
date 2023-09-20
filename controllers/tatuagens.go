package controller

import (
	"fmt"
	"net/http"
)

type Controller struct{}

func (c Controller) Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bem-vindo ao nosso site!")
}
