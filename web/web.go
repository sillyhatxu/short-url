package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"short-url/bolt"
	"short-url/conf"
	"short-url/dto"
	"short-url/short"
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
	router.HandleFunc("/url-list", GetList).Methods(http.MethodGet)
	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	log.Fatal(err)
}

func GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	boltClient := bolt.NewBoltClient(conf.Conf.DBName, 0600)
	result, err := boltClient.ForEach()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}
	var compressArray []dto.CompressDTO
	for _, value := range result {
		var compress dto.CompressDTO
		json.Unmarshal([]byte(value), &compress)
		compressArray = append(compressArray, compress)
	}
	compressArrayJSON, _ := json.Marshal(compressArray)
	w.Write(compressArrayJSON)
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
	}
	if strings.ToLower(longURL.Scheme) != "http" && strings.ToLower(longURL.Scheme) != "https" {
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: "requested url is not a http or https url"})
		w.Write(errMsg)
		return
	}

	var shortenedURL string
	shortenedURL, err = short.Short(shortReq.LongURL)
	if err != nil {
		log.Printf("short url error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}
	shortenedURL = (&url.URL{
		Scheme: conf.Conf.Schema,
		Host:   conf.Conf.DomainName,
		Path:   shortenedURL,
	}).String()
	shortResp, err := json.Marshal(dto.ShortResp{ShortURL: shortenedURL})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}
	w.Write(shortResp)
}

func UncompressURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read expand request error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}

	var expandReq dto.ExpandReq
	err = json.Unmarshal(body, &expandReq)
	if err != nil {
		log.Printf("parse expand request error. %v", err)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}

	var shortURL *url.URL
	shortURL, err = url.Parse(expandReq.ShortURL)
	if err != nil {
		log.Printf(`"%v" is not a valid url`, expandReq.ShortURL)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}
	var expandedURL string
	expandedURL, err = short.Expand(strings.TrimLeft(shortURL.Path, "/"))
	if err != nil {
		log.Printf("expand url error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(dto.ErrorResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}
	expandResp, _ := json.Marshal(dto.ExpandResp{LongURL: expandedURL})
	w.Write(expandResp)
}
func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortededURL := vars["shortenedURL"]

	longURL, err := short.Expand(shortededURL)
	if err != nil {
		log.Printf("redirect short url error. %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	if len(longURL) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
