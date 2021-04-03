package shortener

import (
	"net/http"
)

type RedirectApi interface {
	GetApiHandler() (http.Handler, error)
}
