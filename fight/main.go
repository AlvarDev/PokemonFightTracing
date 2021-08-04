package main

import (
	"encoding/json"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"
	"go.opencensus.io/plugin/ochttp"
)

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/", rootHandler)

	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)
	handler := middleware(mux)

	httpHandler := &ochttp.Handler{
		Propagation: &propagation.HTTPFormat{},
		Handler:     handler,
	}

	log.Info().Msg("Serving pokemon fight")

	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		log.Fatal().Err(err).Msg("Can’t start service")
	}
}

// Response definition for response API
type Response struct {
	Message string `json:"message"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	// re := regexp.MustCompile("/type/([0-9]+)*")
	response := Response{}
	response.Message = "Hello world!"
	responseJSON, err := json.Marshal(response)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse json")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}
