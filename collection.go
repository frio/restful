package restful

import (
	"encoding/json"
	// "github.com/gorilla/mux"
	"net/http"
)

type Collection struct {
	Get  func() ([]interface{}, error)
	Post func(create json.Decoder) (interface{}, error)
}

func (c *Collection) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Add("Content-Type", "application/json")

	if req.Method == "GET" && c.Get != nil {
		value, err := c.Get()

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		encoded, _ := json.Marshal(value)
		rw.Write(encoded)
	} else if req.Method == "POST" && c.Post != nil {
		pleaseCreate := json.NewDecoder(req.Body)

		created, err := c.Post(*pleaseCreate)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		encoded, _ := json.Marshal(created)
		rw.WriteHeader(http.StatusCreated)
		rw.Write(encoded)
	} else {
		allowed := ""

		if c.Get != nil {
			allowed += "GET"
		}

		if c.Post != nil {
			allowed += ", POST"
		}

		rw.Header().Add("Allow", allowed)
		http.Error(rw, "", http.StatusMethodNotAllowed)
	}

}
