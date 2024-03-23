package http

import "net/http"

const (
	swaggeruiprefix = "/swagger-ui/"
	uiDirPath       = "./third_party/swagger-ui"
)

func swaggerUIHandler(swaggerPath string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(swaggeruiprefix+"swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerPath)
	})
	ui := http.FileServer(http.Dir(uiDirPath))
	mux.Handle(swaggeruiprefix, http.StripPrefix(swaggeruiprefix, ui))
	return mux
}
