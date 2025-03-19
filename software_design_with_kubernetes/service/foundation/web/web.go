package web

import (
	"context"
	"net/http"
	"os"

	"github.com/dimfeld/httptreemux"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)   // route specific middleware like authentication
	handler = wrapMiddleware(a.mw, handler) // global middleware like logging

	h := func(w http.ResponseWriter, r *http.Request) {

		// ADD ANY LOGIC

		if err := handler(r.Context(), w, r); err != nil {
			// TODO
			return
		}
	}

	a.ContextMux.Handle(method, path, h)
}
