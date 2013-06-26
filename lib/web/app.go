// The Babou web application core
package web

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

// A `Controller` handles a request by taking an action-name
type Controller interface {
	HandleRequest(string, map[string]string) *Result
}

type DevController interface {
	HandleRequest(string) *Result
	SetContext(Context) error
}

// A route is part of a controller that is capable
// of managing instances for a request life-cycle.
type Route interface {
	Process(string, Context) (DevController, error)
	NewInstance() DevController
	IsSafeInstance() bool // Can this handle requests?
}

// An action takes a map of request-parameters from the middleware
// or router and turns it into a servicable HTTP result.
type Action func(map[string]string) *Result

// Interface for the most basic context: one which encapsulates request parameters.
type Context interface {
	SetParams(map[string]string)
	GetParams() map[string]string
}

// Test impl. of Context interface.
type DevContext struct {
	Params map[string]string
}

func (dc *DevContext) SetParams(params map[string]string) {
	dc.Params = params
}

func (dc *DevContext) GetParams() map[string]string {
	return dc.Params
}

// Represents an HTTP response from a `Controller`
// The middleware or router is responsible for using
// this result appopriately.
type Result struct {
	Body     []byte //HTTP Response Body
	Status   int    //HTTP Status Code
	Redirect *RedirectPath
}

// Requests an HTTP redirect from the middleware or
// router.
type RedirectPath struct {
	NamedRoute string //or:

	ControllerName string
	ActionName     string
}

// Returns a 404 error if a user asks `babou` for the contents of
// a directory. Useful for serving static files.
func DisableDirectoryListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" || strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// Retrieves GET and POST vars from an http Request
func RetrieveAllParams(request *http.Request) map[string]string {
	vars := mux.Vars(request)
	if err := request.ParseForm(); err != nil {
		return vars // could not parse form
	}

	var postVars map[string][]string
	postVars = map[string][]string(request.Form)
	for k, v := range postVars {
		// Ignore duplicate arguments taking the first.
		// POST will supersede any GET data in the event of collisions.
		vars[k] = v[0]
	}

	return vars
}
