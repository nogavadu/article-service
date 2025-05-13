package response

import (
	"github.com/go-chi/render"
	"net/http"
)

type Response struct {
	Error string `json:"error,omitempty"`
}

func Err(w http.ResponseWriter, r *http.Request, errMsg string, status int) {
	render.Status(r, status)
	render.JSON(w, r, &Response{Error: errMsg})
}
