package handlers

import "net/http"

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.FileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}
