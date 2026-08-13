package routes

import (
	"net/http"
	"testing"
)

func TestRegisterRoutesDoesNotPanic(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux, nil)
}
