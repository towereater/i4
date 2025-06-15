package middleware

import (
	"collector/config"
	"net/http"
)

func chainMiddleware(h http.Handler, adapters ...Adapter) http.Handler {
	// for _, a := range slices.Backward(adapters) {
	// 	h = a(h)
	// }
	for i := len(adapters) - 1; i >= 0; i-- {
		h = adapters[i](h)
	}

	return h
}

func AdminAuthentication(h http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adapters := []Adapter{
			Logger(),
			AddConfig(cfg),
			AuthenticateAdmin(),
		}
		chainMiddleware(h, adapters...).ServeHTTP(w, r)
	})
}

func AdminOrFirstAuthentication(h http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adapters := []Adapter{
			Logger(),
			AddConfig(cfg),
			AuthenticateAdminOrFirstUser(),
		}
		chainMiddleware(h, adapters...).ServeHTTP(w, r)
	})
}

func AdminOrClientAuthentication(h http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adapters := []Adapter{
			Logger(),
			AddConfig(cfg),
			AuthenticateAdminOrClient(),
		}
		chainMiddleware(h, adapters...).ServeHTTP(w, r)
	})
}
