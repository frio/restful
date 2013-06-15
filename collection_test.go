package restful_test

import (
	"bytes"
	"encoding/json"
	"github.com/frio/restful"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
)

type CollectionSuite struct{}

var _ = Suite(&CollectionSuite{})

func returnsEmpty() ([]interface{}, error) {
	empty := []interface{}{}
	return empty, nil
}

func createsInstance(decoder json.Decoder) (interface{}, error) {
	var created Example
	err := decoder.Decode(&created)

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *CollectionSuite) TestThatAnEmptyCollectionReturnsAnEmptyList(c *C) {
	request, _ := http.NewRequest("GET", "/test/", nil)

	response := Perform(request, &restful.Collection{
		Get: returnsEmpty,
	})

	c.Assert(response.Code, Equals, 200)
	c.Assert(response.Header().Get("Content-Type"), Equals, "application/json")

	content, err := ioutil.ReadAll(response.Body)
	c.Assert(err, IsNil)

	c.Assert(string(content), Equals, "[]")
}

func (s *CollectionSuite) TestObjectCreation(c *C) {
	expected := &Example{Id: "exists", Changed: true} // set as true because the zero-value of bool is false
	encoded, _ := json.Marshal(expected)

	request, _ := http.NewRequest("POST", "/test/exists/", bytes.NewReader(encoded))
	request.Header.Set("Content-Type", "application/json")

	response := Perform(request, &restful.Collection{
		Post: createsInstance,
	})

	c.Assert(response.Code, Equals, 201)
	c.Assert(response.Header().Get("Content-Type"), Equals, "application/json")

	content, err := ioutil.ReadAll(response.Body)
	c.Assert(err, IsNil)

	var actual Example
	json.Unmarshal(content, &actual)

	c.Assert(actual, Equals, *expected)
}
