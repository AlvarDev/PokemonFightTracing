package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/gorilla/mux"
	"github.com/mtslzr/pokeapi-go"
	"github.com/mtslzr/pokeapi-go/structs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"
	"go.opencensus.io/plugin/ochttp"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	mux := mux.NewRouter()
	mux.HandleFunc("/", rootHandler).Methods("GET", "OPTIONS")

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
	logger := log.Ctx(r.Context())
	logger.Info().Msg("Serving random pokemons")

	ps, err := pokeapi.Resource("pokemon", 1, 500)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	randomSelection := rand.Perm(len(ps.Results))
	re := regexp.MustCompile("/pokemon/([0-9]+)*")

	var result []structs.Pokemon
	for _, v := range randomSelection[:2] {

		match := re.FindStringSubmatch(ps.Results[v].URL)
		p, err := pokeapi.Pokemon(match[1])

		if err != nil {
			log.Ctx(r.Context()).Error().Err(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info().
			Str("pokemonId", match[1]).
			Str("name", p.Name).
			Str("type", p.Types[0].Type.Name).
			Msg("Requesting pokemon")

		result = append(result, p)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.Encode(result)
}
