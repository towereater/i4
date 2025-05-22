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

func LoggedAdminAuthentication(h http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adapters := []Adapter{
			Logger(),
			AddConfig(cfg),
			AuthenticateAdmin(),
		}
		chainMiddleware(h, adapters...).ServeHTTP(w, r)
	})
}

func LoggedAdminFirstAuthentication(h http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adapters := []Adapter{
			Logger(),
			AddConfig(cfg),
			AuthenticateAdminOrFirstUser(),
		}
		chainMiddleware(h, adapters...).ServeHTTP(w, r)
	})
}

func LoggedClientAuthentication(h http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adapters := []Adapter{
			Logger(),
			AddConfig(cfg),
			AuthenticateClient(),
		}
		chainMiddleware(h, adapters...).ServeHTTP(w, r)
	})
}
