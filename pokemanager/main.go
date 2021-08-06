package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"
	"go.opencensus.io/plugin/ochttp"
	"google.golang.org/api/idtoken"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/get-pokemons", pokemonsHandler).Methods("GET")
	r.HandleFunc("/fight-pokemon", fightHandler).Methods("POST")

	rootLogger := zerolog.New(os.Stdout)
	middleware := crzerolog.InjectLogger(&rootLogger)

	handler := cors.Default().Handler(r)
	handler = middleware(handler)

	httpHandler := &ochttp.Handler{
		Propagation: &propagation.HTTPFormat{},
		Handler:     handler,
	}

	log.Info().Msg("Serving pokemon manager")
	if err := http.ListenAndServe(":8080", httpHandler); err != nil {
		log.Fatal().Err(err).Msg("Can’t start service")
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	logger.Info().Msg("Request on get-pokemons")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func pokemonsHandler(w http.ResponseWriter, r *http.Request) {

	// Could be an ENV
	targetURL := "https://pokemon-tbcpnuln2q-rj.a.run.app"

	client, _ := idtoken.NewClient(context.Background(), targetURL)
	log.Ctx(r.Context()).Info().Msg("Request on fight-pokemon")

	req, _ := http.NewRequest("GET", targetURL, nil)
	req = req.WithContext(r.Context())

	resp, err := client.Do(req)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Ctx(r.Context()).Error().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func fightHandler(w http.ResponseWriter, r *http.Request) {

	// Could be an ENV
	targetURL := "https://fight-tbcpnuln2q-rj.a.run.app"

	client, _ := idtoken.NewClient(context.Background(), targetURL)
	log.Ctx(r.Context()).Info().Msg("Request on fight-pokemon")

	req, _ := http.NewRequest("POST", targetURL, r.Body)
	req = req.WithContext(r.Context())

	resp, err := client.Do(req)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Ctx(r.Context()).Error().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
