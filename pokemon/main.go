package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"
	"go.opencensus.io/plugin/ochttp"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	mux := mux.NewRouter()
	mux.HandleFunc("/", rootHandler)

	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)
	handler := middleware(mux)

	httpHandler := &ochttp.Handler{
		Propagation: &propagation.HTTPFormat{},
		Handler:     handler,
	}

	log.Info().Msg("Serving random pokemon")

	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		log.Fatal().Err(err).Msg("Can’t start service")
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{}
	response.Message = "Hello World!"
	responseJSON, err := json.Marshal(response)

	if err != nil {
		log.Fatal().Err(err).Msg("Can’t parse JSON")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
