package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JsonError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{error: %s}`, error)
}

type ErrorResponse struct {
	Error string `json:error`
	Code  int    `json:code`
}

type Response struct {
	Text string `json:"text"`
	Lang string `json:"lang"`
}

func sanitizeHandle(w http.ResponseWriter, r *http.Request) {
	// TODO: add lang
	//lang := r.FormValue("lang")
	//if err != nil {
	//	JsonError(w, "Invalid lang", 400)
	//	return
	//}

	text := r.FormValue("text")
	sanitized := pfilter.Sanitize(text)
	//Logf("lang: %s, text: %s, sanitized: %s", lang, text, sanitized)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(&Response{Text: sanitized})
}

func updateBlacklistHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	blacklist, ok := r.Form["blacklist"]

	if !ok || len(blacklist) == 0 {
		JsonError(w, "Expected `blacklist` key", 400)
		return
	}

	switch r.Method {
		case "PUT":
			pfilter.Update(blacklist)
			w.WriteHeader(200)
		case "POST":
			pfilter.Replace(blacklist)
			w.WriteHeader(201)
		default:
			panic("should not reach")
	}
}

func removeBlacklistHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	blacklist, ok := r.Form["blacklist"]

	if !ok || len(blacklist) == 0 {
		JsonError(w, "Expected `blacklist` key", 400)
		return
	}

	pfilter.Remove(blacklist)
	w.WriteHeader(200)
}

func getBlacklistHandle(w http.ResponseWriter, r *http.Request) {
	blacklist := pfilter.Blacklist()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(&blacklist)
}
