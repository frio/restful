package restful_test

import (
	"github.com/gorilla/mux"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

func Perform(request *http.Request, handler http.Handler) *httptest.ResponseRecorder {
	router := mux.NewRouter()
	router.Handle("/test/", handler)
	router.Handle("/test/{id}/", handler)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	return w
}

type Example struct {
	Id      string
	Changed bool
}
