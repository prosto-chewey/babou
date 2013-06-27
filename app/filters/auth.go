package filters

import (
	"errors"
	"fmt"
	"net/http"

	web "babou/lib/web"
)

// An authorizable route must be AuthContext aware.
// That way it gets access to all helpers defined by the context.
type AuthorizableRoute interface {
	Process(string) (AuthorizableController, error)
	NewInstance() AuthorizableController
}

type AuthorizableController interface {
	web.Controller
	SetAuthContext(*AuthContext) error
}

// An impl. of SessionContext that uses it to provide helper methods for auth'ing a user.
type AuthContext struct {
	params map[string]string
	isInit bool
}

// Returns an uninitialized AuthContext suitable for use in a context-chain
func AuthChain() *AuthContext {
	context := &AuthContext{isInit: false}

	return context
}

// Tests if the route implements AuthorizableController interface as well as the SessionChain interface.
func (ac *AuthContext) TestContext(route web.Route, chain []web.ChainableContext) error {
	//requires AuthorizableController and SessionChain
	hasSession := false

	_, ok := route.(AuthorizableController)
	if !ok {
		return errors.New(fmt.Sprintf("The route :: %T :: does not support the AuthContext.", route))
	}

	for i := 0; i < len(chain); i++ {
		_, ok := chain[i].(SessionChainLink)
		if ok {
			hasSession = true
		}
	}

	if hasSession && ok {
		return nil
	} else {
		return errors.New(fmt.Sprintf("The route :: %T :: does not have a SessionAware context in it's context chain.", route))
	}
}

func (ac *AuthContext) NewInstance() web.ChainableContext {
	newAc := &AuthContext{isInit: false}

	return newAc
}

// Implements ChainableContext
func (ac *AuthContext) ApplyContext(controller web.Controller, response http.ResponseWriter, request *http.Request, chain []web.ChainableContext) {
	ac.isInit = true

	v, ok := controller.(AuthorizableController)
	if ok {
		if err := v.SetAuthContext(ac); err != nil {
			fmt.Printf("Error setting authorization context: %s \n", err.Error())
		}
	} else {
		fmt.Printf("Tried to wrap a controller that is not AuthContext aware \n")
	}
}

// Returns `true` if the AuthContext is properly initialized on top of a session store.
func (ac *AuthContext) isValid() bool {
	return ac.isInit
}
