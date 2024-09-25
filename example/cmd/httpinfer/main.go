package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tarcisio/infergo"
	"github.com/tarcisio/infergo/example"
)

type Service struct {
	engine *infergo.Engine[*example.Payload]
}

func NewService(engine *infergo.Engine[*example.Payload]) *Service {
	return &Service{engine: engine}
}

func (service *Service) ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	var payload example.Payload

	// 1. Decode the input payload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// 2. Execute the rule engine
	if err := service.engine.Execute(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Return the modified payload as JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	eng := infergo.NewEngine[*example.Payload](100)
	eng.AddRule(example.AgeRule(), 10)
	eng.AddRule(example.StateRule(), 20)

	svc := NewService(eng)

	http.HandleFunc("POST /execute", svc.ExecuteHandler)

	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Fatal(err)
	}
}
