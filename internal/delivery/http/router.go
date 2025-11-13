package http

import (
	"net/http"

	"github.com/valentinfrappart/securerestapi/internal/infrastructure/security"
)

type Router struct {
	handler    *Handler
	jwtService *security.JWTService
}

func NewRouter(handler *Handler, jwtService *security.JWTService) *Router {
	return &Router{
		handler:    handler,
		jwtService: jwtService,
	}
}

func (rt *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", applyMiddlewares(rt.handler.Health, CORSMiddleware, LoggingMiddleware))
	mux.HandleFunc("/api/auth/register", applyMiddlewares(rt.handler.Register, CORSMiddleware, LoggingMiddleware))
	mux.HandleFunc("/api/auth/login", applyMiddlewares(rt.handler.Login, CORSMiddleware, LoggingMiddleware))

	mux.HandleFunc("/api/auth/me", applyMiddlewares(rt.handler.Me, CORSMiddleware, LoggingMiddleware, AuthMiddleware(rt.jwtService)))

	return mux
}

func applyMiddlewares(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
