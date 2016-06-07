package ui

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"sourcegraph.com/sourcegraph/sourcegraph/app/internal"
	"sourcegraph.com/sourcegraph/sourcegraph/app/internal/tmpl"
	approuter "sourcegraph.com/sourcegraph/sourcegraph/app/router"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/errcode"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/handlerutil"
)

func init() {
	router.Get(routeBlob).Handler(handler(serveBlob))
	router.Get(routeBuild).Handler(handler(serveBuild))
	router.Get(routeDef).Handler(handler(serveDef))
	router.Get(routeDefInfo).Handler(handler(serveDefInfo))
	router.Get(routeRepo).Handler(handler(serveRepo))
	router.Get(routeRepoBuilds).Handler(handler(serveRepoBuilds))
	router.Get(routeTree).Handler(handler(serveTree))
	router.Get(routeTopLevel).Handler(internal.Handler(serveAny))
	router.PathPrefix("/").Methods("GET").Handler(internal.Handler(serveAny))
	router.NotFoundHandler = internal.Handler(serveAny)

	// Attach to app handler's catch-all UI route. This is better than
	// adding the UI routes to the app router directly because it
	// keeps the two routers separate.
	internal.Handlers[approuter.UI] = router
}

// handler wraps h, calling tmplExec with the HTTP equivalent error
// code of h's return value (or HTTP 200 if err == nil).
func handler(h func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return internal.Handler(func(w http.ResponseWriter, r *http.Request) error {
		err := h(w, r)
		return tmplExec(w, r, errcode.HTTP(err))
	})
}

// These handlers return the proper HTTP status code but otherwise do
// not pass data to the JavaScript UI code.

func serveBlob(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := handlerutil.Client(r)
	_, err := handlerutil.GetRepo(ctx, mux.Vars(r))
	if err != nil {
		return err
	}
	return nil
}

func serveBuild(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := handlerutil.Client(r)
	_, err := handlerutil.GetRepo(ctx, mux.Vars(r))
	if err != nil {
		return err
	}
	return nil
}

func serveDef(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := handlerutil.Client(r)
	_, err := handlerutil.GetRepo(ctx, mux.Vars(r))
	if err != nil {
		return err
	}
	return nil
}

func serveDefInfo(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := handlerutil.Client(r)
	_, err := handlerutil.GetRepo(ctx, mux.Vars(r))
	if err != nil {
		return err
	}
	return nil
}

func serveRepo(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := handlerutil.Client(r)
	_, err := handlerutil.GetRepo(ctx, mux.Vars(r))
	if err != nil {
		return err
	}
	return nil
}

func serveRepoBuilds(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := handlerutil.Client(r)
	_, err := handlerutil.GetRepo(ctx, mux.Vars(r))
	if err != nil {
		return err
	}
	return nil
}

func serveTree(w http.ResponseWriter, r *http.Request) error {
	ctx, _ := handlerutil.Client(r)
	_, err := handlerutil.GetRepo(ctx, mux.Vars(r))
	if err != nil {
		return err
	}
	return nil
}

// serveAny is the fallback/catch-all route. It preloads nothing and
// returns a page that will merely bootstrap the JavaScript app.
func serveAny(w http.ResponseWriter, r *http.Request) error {
	return tmplExec(w, r, http.StatusOK)
}

func tmplExec(w http.ResponseWriter, r *http.Request, statusCode int) error {
	return tmpl.Exec(r, w, "ui.html", statusCode, nil, &struct {
		tmpl.Common
		Stores *json.RawMessage
	}{})
}
