package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"short-url/conf"
	"short-url/dto"
	"strings"
	"time"
)

func Start() {
	log.Info("web starts")
	router := mux.NewRouter()

	router.HandleFunc("/health", CheckHealth).Methods(http.MethodGet)
	router.HandleFunc("/compress-url", CompressURL).Methods(http.MethodPost).HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/uncompress-url", UncompressURL).Methods(http.MethodPost).HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/{shortenedURL:[a-zA-Z0-9]{1,11}}", Redirect).Methods(http.MethodGet)
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	log.Fatal(err)
}

func CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	health, _ := json.Marshal(dto.Health{Status: "UP", Description: "Short-URL Success"})
	w.Write(health)
}

func CompressURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read short request error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}

	var shortReq dto.ShortReq
	err = json.Unmarshal(body, &shortReq)
	if err != nil {
		log.Printf("parse short request error. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}

	var longURL *url.URL
	longURL, err = url.Parse(shortReq.LongURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: "requested url is malformed"})
		w.Write(errMsg)
		return
	} else {
		if strings.ToLower(longURL.Scheme) != "http" && strings.ToLower(longURL.Scheme) != "https" {
			w.WriteHeader(http.StatusBadRequest)
			errMsg, _ := json.Marshal(dto.ErrorResp{Msg: "requested url is not a http or https url"})
			w.Write(errMsg)
			return
		}
	}

	var shortenedURL string
	shortenedURL, err = short.Shorter.Short(shortReq.LongURL)
	shortenedURL = (&url.URL{
		Scheme: conf.Conf.Schema,
		Host:   conf.Conf.DomainName,
		Path:   shortenedURL,
	}).String()
	if err != nil {
		log.Printf("short url error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(errorResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	} else {
		shortResp, _ := json.Marshal(shortResp{ShortURL: shortenedURL})
		w.Write(shortResp)
	}
}

func UncompressURL(w http.ResponseWriter, r *http.Request) {

}
func Redirect(w http.ResponseWriter, r *http.Request) {

}
