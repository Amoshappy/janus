package router

import (
	"context"
	"net/http"

	"github.com/dimfeld/httptreemux"
	"github.com/justinas/alice"
)

// ParamsContextKey is used to retrieve a path's params map from a request's context.
const ParamsContextKey = "params.context.key"

type Constructor func(http.Handler) http.Handler

type Params struct {
	params map[string]string
}

// FromContext returns the params map associated with the given context if one exists. Otherwise, an empty map is returned.
func FromContext(ctx context.Context) *Params {
	if p, ok := ctx.Value(ParamsContextKey).(map[string]string); ok {
		return &Params{p}
	}
	return &Params{}
}

func (p *Params) ByName(name string) string {
	return p.params[name]
}

type Router interface {
	Handle(method string, path string, handler http.HandlerFunc, handlers ...Constructor)
	Any(path string, handler http.HandlerFunc, handlers ...Constructor)
	GET(path string, handler http.HandlerFunc, handlers ...Constructor)
	POST(path string, handler http.HandlerFunc, handlers ...Constructor)
	PUT(path string, handler http.HandlerFunc, handlers ...Constructor)
	DELETE(path string, handler http.HandlerFunc, handlers ...Constructor)
	PATCH(path string, handler http.HandlerFunc, handlers ...Constructor)
	HEAD(path string, handler http.HandlerFunc, handlers ...Constructor)
	OPTIONS(path string, handler http.HandlerFunc, handlers ...Constructor)
	Group(path string) Router
	Use(handlers ...Constructor) Router
}

type HttpTreeMuxRouter struct {
	mux         *httptreemux.TreeMux
	innerRouter *httptreemux.ContextGroup
	chain       alice.Chain
}

type Options struct {
	NotFoundHandler           http.HandlerFunc
	SafeAddRoutesWhileRunning bool
	RedirectMethodBehavior    map[string]httptreemux.RedirectBehavior
}

var DefaultOptions = Options{
	NotFoundHandler:           http.NotFound,
	SafeAddRoutesWhileRunning: true,
	RedirectMethodBehavior: map[string]httptreemux.RedirectBehavior{
		// tree mux router uses Redirect301 behavior by default for paths that differs with slash at the end
		// from registered, that causes problems with some services, e.g. api-v2 OPTIONS /menus/ gives 301 and we
		// want by-pass it as is, but only for OPTIONS method, that is processed by CORS
		http.MethodOptions: httptreemux.UseHandler,
	},
}

func NewHttpTreeMuxWithOptions(options Options) *HttpTreeMuxRouter {
	router := httptreemux.New()
	router.NotFoundHandler = options.NotFoundHandler
	router.SafeAddRoutesWhileRunning = options.SafeAddRoutesWhileRunning
	router.RedirectMethodBehavior = options.RedirectMethodBehavior

	return &HttpTreeMuxRouter{
		mux:         router,
		innerRouter: router.UsingContext(),
		chain:       alice.New(),
	}
}

func NewHttpTreeMuxRouter() *HttpTreeMuxRouter {
	return NewHttpTreeMuxWithOptions(DefaultOptions)
}

func (r *HttpTreeMuxRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *HttpTreeMuxRouter) Any(path string, handler http.HandlerFunc, handlers ...Constructor) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodHead,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, method := range methods {
		r.Handle(method, path, handler, handlers...)
	}
}

func (r *HttpTreeMuxRouter) Handle(method string, path string, handler http.HandlerFunc, handlers ...Constructor) {
	var chain alice.Chain
	chain = r.chain

	for _, h := range handlers {
		chain = chain.Append(alice.Constructor(h))
	}

	r.innerRouter.Handle(method, path, chain.ThenFunc(handler).ServeHTTP)
}

func (r *HttpTreeMuxRouter) GET(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.Handle(http.MethodGet, path, handler, handlers...)
}

func (r *HttpTreeMuxRouter) POST(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.Handle(http.MethodPost, path, handler, handlers...)
}

func (r *HttpTreeMuxRouter) PUT(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.Handle(http.MethodPut, path, handler, handlers...)
}

func (r *HttpTreeMuxRouter) DELETE(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.Handle(http.MethodDelete, path, handler, handlers...)
}

func (r *HttpTreeMuxRouter) PATCH(path string, handler http.HandlerFunc, handlers ...Constructor) {
	r.Handle(http.MethodPatch, path, handler, handlers...)
}

func (r *HttpTreeMuxRouter) HEAD(path string, handler http.HandlerFunc, middleware ...Constructor) {
	r.Handle(http.MethodHead, path, handler, middleware...)
}

func (r *HttpTreeMuxRouter) OPTIONS(path string, handler http.HandlerFunc, middleware ...Constructor) {
	r.Handle(http.MethodOptions, path, handler, middleware...)
}

func (r *HttpTreeMuxRouter) Group(path string) Router {
	return &HttpTreeMuxRouter{
		mux:         r.mux,
		innerRouter: r.innerRouter.NewContextGroup(path),
		chain:       r.chain,
	}
}

func (r *HttpTreeMuxRouter) Use(handlers ...Constructor) Router {
	for _, h := range handlers {
		r.chain = r.chain.Append(alice.Constructor(h))
	}

	return r
}

func (r *HttpTreeMuxRouter) wrapConstructor(handlers []Constructor) []alice.Constructor {
	var cons = make([]alice.Constructor, len(handlers))
	for _, m := range handlers {
		cons = append(cons, alice.Constructor(m))
	}
	return cons
}
