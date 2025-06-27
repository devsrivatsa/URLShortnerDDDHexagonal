package api

import (
	"fmt"
	"io"
	"log"
	"net/http"

	js "github.com/devsrivatsa/URLShortnerDDDHexagonal/serializer/json"
	mp "github.com/devsrivatsa/URLShortnerDDDHexagonal/serializer/msgpack"
	"github.com/devsrivatsa/URLShortnerDDDHexagonal/urlShortner"
	chi "github.com/go-chi/chi/v5"
	errs "github.com/pkg/errors"
)

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

type handler struct {
	redirectService urlShortner.RedirectService
}

func NewRedirectHandler(rs urlShortner.RedirectService) RedirectHandler {
	return &handler{
		redirectService: rs,
	}
}

func setupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (h *handler) serializer(contentType string) urlShortner.RedirectSerializer {
	if contentType == "application/json" {
		return &js.Redirect{}
	}
	return &mp.Redirect{}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "code parsed from url is empty", http.StatusBadRequest)
		return
	}
	log.Printf("Getting redirect for code: %s", code)
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errs.Cause(err) == urlShortner.ErrRedirectNotFound {
			http.Error(w, fmt.Sprintf("redirect not found: %v", err), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	redirect, err := h.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = h.redirectService.Store(redirect)
	if err != nil {
		if errs.Cause(err) == urlShortner.ErrRedirectInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseBody, err := h.serializer(contentType).Encode(redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(w, contentType, responseBody, http.StatusCreated)
}
