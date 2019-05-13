package routes

import (
	"net/http"

	"github.com/go-chi/render"
)

type Error struct {
	Code         int    `json:"code"`
	Description  string `json:"description"`
	ReasonPhrase string `json:"reasonPhrase"`
}

func RenderError(w http.ResponseWriter, r *http.Request, message string, statusCode int) {
	reasonPhrase := ""
	switch statusCode {
	case http.StatusBadRequest:
		reasonPhrase = "Bad Request"
	case http.StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	case http.StatusNotFound:
		reasonPhrase = "Not Found"
	case http.StatusConflict:
		reasonPhrase = "Conflict"
	default:
		reasonPhrase = "Unknown error"
	}

	render.Status(r, statusCode)
	render.JSON(w, r, &Error{
		Code:         statusCode,
		Description:  message,
		ReasonPhrase: reasonPhrase,
	})
}
