package main

import (
	// "io"
	"log"
	"fmt"
	"os"
	"net/http"
	// "time"
	"encoding/json"
	"github.com/gorilla/mux"
    "github.com/gorilla/handlers"
)

const (
	maxFish = 5
)

type myHandler struct{}

type Submarine struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Artifact struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Fish struct {
	X int `json:"x"`
	Y int `json:"y"`
}

var submarine Submarine
var artifact Artifact
var fish []Fish

func handleMessage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    message := vars["msg"]
    response := map[string]string{"message": message}
    json.NewEncoder(w).Encode(response)
}

func handleNumber(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    number := vars["num"]
    response := map[string]string{"number": number}
    json.NewEncoder(w).Encode(response)
}

func handleGetSubmarine(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(submarine)
}

func handleGetArtifact(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(artifact)
}

func handleGetFish(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(fish)
}


func handleMoveSubmarine(w http.ResponseWriter, r *http.Request) {
	
	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json\n"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var newSubmarine Submarine
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newSubmarine)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%+v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		submarine.X += newSubmarine.X
		submarine.Y += newSubmarine.Y
		fmt.Fprintf(os.Stdout, "Submarine moved to: %+v\n", submarine)
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(submarine)
}

func handleUpdateArtifact(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json\n"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var newArtifact Artifact
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newArtifact)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%+v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		artifact.X = newArtifact.X
		artifact.Y = newArtifact.Y
		fmt.Fprintf(os.Stdout, "Artifact placed at: %+v\n", artifact)
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(artifact)
}

func handleAddFish(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		msg := "Content-Type header is not application/json\n"
		http.Error(w, msg, http.StatusUnsupportedMediaType)
		return
	}

	var newFish Fish
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&newFish)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%+v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		if len(fish) < maxFish {
			fish = append(fish, newFish)
			fmt.Fprintf(os.Stdout, "New Fish: %+v\n", newFish)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newFish)

		} else {
			msg := "Maximum number of fish excedeed\n"
			fmt.Fprintf(os.Stdout, "%s\n", msg)
			http.Error(w, msg, http.StatusBadRequest)
		}
	}
}


func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}


func main() {

	port := ":8080"
	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST"})

    var router = mux.NewRouter()
	router.Use(commonMiddleware)

	fish = make([]Fish, 0)
	submarine = Submarine {
		X: 10,
		Y: 15,
	}
	
    router.HandleFunc("/api/submarine", handleGetSubmarine).Methods("GET")
    router.HandleFunc("/api/submarine/move", handleMoveSubmarine).Methods("POST")
    router.HandleFunc("/api/artifact", handleGetArtifact).Methods("GET")
    router.HandleFunc("/api/artifact/update", handleUpdateArtifact).Methods("POST")
    router.HandleFunc("/api/fish", handleGetFish).Methods("GET")
    router.HandleFunc("/api/fish/add", handleAddFish).Methods("POST")

    fmt.Printf("Server is running at http://localhost%s\n", port)
    log.Fatal(http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
