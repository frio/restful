package restful

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Resource struct {
	Get    func(id string) (interface{}, error)
	Delete func(id string) error
	Put    func(from interface{}, to json.Decoder) (interface{}, error)
	// no POST because POST-ing to a resource doesn't make sense
}

func (r *Resource) String() string {
	return "test"
}

func (r *Resource) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	rw.Header().Add("Content-Type", "application/json")

	id, ok := vars["id"]

	if !ok {
		http.Error(rw, "No ID -- Should never happen!", 500)
	}

	if req.Method == "GET" && r.Get != nil {
		value, _ := r.Get(id)

		if value == nil {
			http.NotFound(rw, req)
		}

		encoded, _ := json.Marshal(value)
		rw.Write(encoded)
	} else if req.Method == "DELETE" && r.Delete != nil {
		err := r.Delete(id)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	} else if (req.Method == "PUT" || req.Method == "PATCH") && r.Get != nil && r.Put != nil {
		from, _ := r.Get(id)

		if from == nil {
			http.NotFound(rw, req)
			return
		}

		to := json.NewDecoder(req.Body)

		final, err := r.Put(from, *to)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		encoded, _ := json.Marshal(final)
		rw.Write(encoded)
	} else {
		allowed := ""

		if r.Get != nil {
			allowed += "GET"
		}

		if r.Delete != nil {
			allowed += ", DELETE"
		}

		if r.Put != nil {
			allowed += ", PUT"
		}

		rw.Header().Add("Allow", allowed)
		http.Error(rw, "", http.StatusMethodNotAllowed)
	}
}
