package main

import (
	// "io"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	maxFish        = 5
	httpServerPort = ":8080"
	tcpServerPort  = ":8000"
)

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
var tcpConn net.Conn

func toJson(v interface{}) string {
	msg, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(msg)
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

		sendData(tcpConn, "UpdateSubmarine:"+toJson(submarine))
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

		sendData(tcpConn, "UpdateArtifact:"+toJson(artifact))
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

			sendData(tcpConn, "Newfish:"+toJson(newFish))

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

func initTcpServer() {
	fmt.Println("TCP server started...")
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("Error starting socket server: " + err.Error())
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error listening to client: " + err.Error())
	}
	tcpConn = conn
	fmt.Println(conn.RemoteAddr().String() + ": client connected")
}

func sendData(conn net.Conn, data string) {
	_, err := fmt.Fprintf(conn, data+"\n")
	if err != nil {
		fmt.Println(conn.RemoteAddr().String() + ": end sending data")
		return
	}
}

func main() {

	// wg := new(sync.WaitGroup)
	// wg.Add(2)

	headersOk := handlers.AllowedHeaders([]string{"Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST"})

	var router = mux.NewRouter()
	router.Use(commonMiddleware)

	fish = make([]Fish, 0)
	submarine = Submarine{
		X: 10,
		Y: 15,
	}

	router.HandleFunc("/api/submarine", handleGetSubmarine).Methods("GET")
	router.HandleFunc("/api/submarine/move", handleMoveSubmarine).Methods("POST")
	router.HandleFunc("/api/artifact", handleGetArtifact).Methods("GET")
	router.HandleFunc("/api/artifact/update", handleUpdateArtifact).Methods("POST")
	router.HandleFunc("/api/fish", handleGetFish).Methods("GET")
	router.HandleFunc("/api/fish/add", handleAddFish).Methods("POST")

	// tcpConn = initTcpServer()
	// if tcpConn == nil {
	// 	return
	// }

	go initTcpServer()

	fmt.Printf("Server is running at http://localhost%s\n", httpServerPort)
	http.ListenAndServe(httpServerPort, handlers.CORS(originsOk, headersOk, methodsOk)(router))
}
