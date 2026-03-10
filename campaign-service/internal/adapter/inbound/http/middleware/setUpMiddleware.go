package middleware

import (
	"net/http"

	"github.com/likhon22/ad-system/campaign-service/internal/utils"
)

func SetUpMiddleware(mux http.Handler, mw *Middleware) http.Handler {

	return utils.ChainMiddleware(
		mux,
		mw.Logger,
	)

}
