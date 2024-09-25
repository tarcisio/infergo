package main

import (
	"encoding/json"
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

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Execute the rule engine
	err = service.engine.Execute(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the modified payload as JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func main() {
	eng := infergo.NewEngine[*example.Payload](100)
	eng.AddRule(example.AgeRule(), 10)
	eng.AddRule(example.StateRule(), 20)

	svc := NewService(eng)

	http.HandleFunc("POST /execute", svc.ExecuteHandler)
	http.ListenAndServe(":8181", nil)
}
