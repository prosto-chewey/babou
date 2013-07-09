package app

import (
	http "net/http"

	controllers "github.com/drbawb/babou/app/controllers"
	filters "github.com/drbawb/babou/app/filters"

	web "github.com/drbawb/babou/lib/web"

	mux "github.com/gorilla/mux"
)

func LoadRoutes() *mux.Router {
	r := mux.NewRouter()
	web.Router = r

	// Shorthand for controllers
	home := controllers.NewHomeController()
	login := controllers.NewLoginController()
	torrent := controllers.NewTorrentController()

	// Shows public homepage, redirects to private site if valid session can be found.
	r.HandleFunc("/",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(home, "index")).Name("homeIndex")

	// Displays a login form.
	r.HandleFunc("/login",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "index")).
		Methods("GET").
		Name("loginIndex")

	r.HandleFunc("/login",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "session")).
		Methods("POST").
		Name("loginSession")

	r.HandleFunc("/logout",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "logout")).
		Methods("GET").
		Name("loginDelete")

	// Displays a registration form
	r.HandleFunc("/register",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "new")).Methods("GET").Name("loginNew")
	// Handles a new user's registration request.
	r.HandleFunc("/register",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(login, "create")).Methods("POST").Name("loginCreate")

	// Handle torrent routes:
	r.HandleFunc("/torrents",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(torrent, "index")).Methods("GET").Name("torrentIndex")

	r.HandleFunc("/torrents/new",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(torrent, "new")).Methods("GET").Name("torrentNew")

	r.HandleFunc("/torrents/create",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(torrent, "create")).Methods("POST").Name("torrentCreate")
	r.HandleFunc("/torrents/download/{torrentId}",
		filters.BuildDefaultChain().
			Chain(filters.AuthChain()).
			Execute(torrent, "download")).
		Methods("GET").
		Name("torrentDownload")

	// Catch-All: Displays all public assets.
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/",
		web.DisableDirectoryListing(http.FileServer(http.Dir("assets/")))))

	return r
}

//TODO: should move some of this to a library package.
func handleRedirect(redirect *web.RedirectPath, response http.ResponseWriter, request *http.Request) {
	if redirect.NamedRoute != "" {
		url, err := web.Router.Get(redirect.NamedRoute).URL()
		if err != nil {
			http.Error(response, string("While trying to redirect you to another page the server encountered an error. Please reload the homepage"),
				500)
		}

		http.Redirect(response, request, url.Path, 302)
	}
}
