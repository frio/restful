package restful_test

import (
	"bytes"
	"encoding/json"
	"github.com/frio/restful"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"net/http"
)

type ResourceSuite struct{}

var _ = Suite(&ResourceSuite{})

func simpleUpdate(from interface{}, request json.Decoder) (interface{}, error) {
	assertedFrom := from.(*Example)
	var to Example
	err := request.Decode(&to)
	if err != nil {
		return nil, err
	}
	return &Example{Id: assertedFrom.Id, Changed: to.Changed}, nil
}

func returnsNothing(id string) (interface{}, error) {
	return nil, nil
}

func returnsSomething(id string) (interface{}, error) {
	return &Example{Id: id, Changed: false}, nil
}

func (s *ResourceSuite) TestThatAnUnkownResourceReturns404(c *C) {
	request, _ := http.NewRequest("GET", "/test/nonexistant", nil)

	response := Perform(request, &restful.Resource{
		Get: returnsNothing,
	})

	c.Assert(response.Code, Equals, 404)
}

func (s *ResourceSuite) TestThatAKnownResourceReturnsAJSONEncoding(c *C) {
	request, _ := http.NewRequest("GET", "/test/exists/", nil)

	response := Perform(request, &restful.Resource{
		Get: returnsSomething,
	})

	c.Assert(response.Code, Equals, 200)
	c.Assert(response.Header().Get("Content-Type"), Equals, "application/json")

	content, err := ioutil.ReadAll(response.Body)
	c.Assert(err, IsNil)

	var body map[string]string
	json.Unmarshal(content, &body)

	id, _ := body["Id"]

	c.Assert(id, Equals, "exists")
}

func (s *ResourceSuite) TestThatGetToldAboutUnsupportedMethods(c *C) {
	request, _ := http.NewRequest("DUMB", "/test/irrelevant/", nil)

	response := Perform(request, &restful.Resource{
		Get: returnsNothing,
		Delete: func(id string) error {
			return nil
		},
	})

	c.Assert(response.Code, Equals, http.StatusMethodNotAllowed)
	c.Assert(response.Header().Get("Allow"), Equals, "GET, DELETE")
}

func (s *ResourceSuite) TestThatWeCanDeleteAResource(c *C) {
	request, _ := http.NewRequest("DELETE", "/test/exists/", nil)

	response := Perform(request, &restful.Resource{
		Get: returnsSomething,
		Delete: func(id string) error {
			return nil
		},
	})

	c.Assert(response.Code, Equals, 200)
	c.Assert(response.Header().Get("Content-Type"), Equals, "application/json")
}

func (s *ResourceSuite) TestThatWeCanUpdateAResource(c *C) {
	expected := &Example{Id: "exists", Changed: true}
	encoded, _ := json.Marshal(expected)
	request, _ := http.NewRequest("PATCH", "/test/exists/", bytes.NewReader(encoded))
	request.Header.Set("Content-Type", "application/json")

	response := Perform(request, &restful.Resource{
		Get: returnsSomething,
		Put: simpleUpdate,
	})

	c.Assert(response.Code, Equals, 200)
	c.Assert(response.Header().Get("Content-Type"), Equals, "application/json")

	content, err := ioutil.ReadAll(response.Body)
	c.Assert(err, IsNil)

	var actual Example
	json.Unmarshal(content, &actual)

	c.Assert(actual.Changed, Equals, expected.Changed)
}
