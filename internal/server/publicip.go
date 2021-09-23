package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/publicip"
)

func newPublicIPHandler(looper publicip.Looper, w warner) http.Handler {
	return &publicIPHandler{
		looper: looper,
		warner: w,
	}
}

type publicIPHandler struct {
	looper publicip.Looper
	warner warner
}

func (h *publicIPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.RequestURI = strings.TrimPrefix(r.RequestURI, "/publicip")
	switch r.RequestURI {
	case "/ip":
		switch r.Method {
		case http.MethodGet:
			h.getPublicIP(w)
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	default:
		http.Error(w, "", http.StatusNotFound)
	}
}

func (h *publicIPHandler) getPublicIP(w http.ResponseWriter) {
	data := h.looper.GetData()
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(data); err != nil {
		h.warner.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
