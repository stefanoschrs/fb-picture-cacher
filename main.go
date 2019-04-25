package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

var fbIdRegexp *regexp.Regexp

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fbId := strings.Split(req.URL.Path, "/")[1]
	if !fbIdRegexp.MatchString(fbId) {
		handleError(w, http.StatusBadRequest)
		return
	}

	res, err := http.Get("http://graph.facebook.com/" + fbId + "/picture?type=large")
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Println(err)
		handleError(w, http.StatusInternalServerError)
		return
	}
}

func handleError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	_, err := w.Write([]byte(http.StatusText(status)))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	port := ":8080"
	if env, ok := os.LookupEnv("PORT"); ok {
		port = ":" + env
	}

	fbIdRegexp = regexp.MustCompile(`^[0-9]+$`)

	// Cache
	memcached, _ := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)

	cacheClient, _ := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(7 * 24 * time.Hour),
		cache.ClientWithRefreshKey("opn"),
	)

	// Handlers
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		handleError(w, http.StatusNotFound)
		return
	})
	http.Handle("/", cacheClient.Middleware(http.HandlerFunc(indexHandler)))

	log.Printf("Listening on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
