package utils

import "net/http"

func ChainMiddleware(mux http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		mux = mws[i](mux)

	}
	return mux
}
