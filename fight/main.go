package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/gorilla/mux"
	"github.com/mtslzr/pokeapi-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yfuruyama/crzerolog"
	"go.opencensus.io/plugin/ochttp"
)

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/", rootHandler).Methods("POST")

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

// Request struct handler
type PokemonsFight struct {
	PokemonA struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Photo string `json:"photo"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"pokemonA"`

	PokemonB struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Photo string `json:"photo"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"pokemonB"`
}

// Response definition for response API
type Response struct {
	Message string `json:"message"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	var pf PokemonsFight
	err := json.NewDecoder(r.Body).Decode(&pf)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	re := regexp.MustCompile("/type/([0-9]+)*")

	matchA := re.FindStringSubmatch(pf.PokemonA.Types[0].Type.URL)
	matchB := re.FindStringSubmatch(pf.PokemonB.Types[0].Type.URL)

	typeA, errA := pokeapi.Type(matchA[1])
	typeB, errB := pokeapi.Type(matchB[1])

	if errA != nil {
		log.Ctx(r.Context()).Error().Err(errA)
		http.Error(w, errA.Error(), http.StatusInternalServerError)
		return
	}

	if errB != nil {
		log.Ctx(r.Context()).Error().Err(errB)
		http.Error(w, errB.Error(), http.StatusInternalServerError)
		return
	}

	damagedTypeA := 0
	damagedTypeB := 0

	for _, v := range typeA.DamageRelations.DoubleDamageFrom {
		if v.Name == pf.PokemonB.Types[0].Type.Name {
			damagedTypeA += 2
		}

		fmt.Print("\t" + v.Name + ": ")
		fmt.Println(damagedTypeA)
	}

	fmt.Print("Total damage taken: ")
	fmt.Println(damagedTypeA)
	fmt.Println("------------")

	fmt.Println("Score B: ")

	for _, v := range typeB.DamageRelations.DoubleDamageFrom {
		if v.Name == pf.PokemonA.Types[0].Type.Name {
			damagedTypeB += 2
		}
		fmt.Print("\t" + v.Name + ": ")
		fmt.Println(damagedTypeB)
	}

	fmt.Print("Total damage taken: ")
	fmt.Println(damagedTypeB)

	if damagedTypeA < damagedTypeB {
		fmt.Println("A won")
	} else if damagedTypeB < damagedTypeA {
		fmt.Println("B won")
	} else {
		fmt.Println("Draw!")
	}

	response := Response{}
	response.Message = "Hello world!"
	responseJSON, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}
