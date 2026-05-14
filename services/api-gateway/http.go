package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	grpc_clients "ride-sharing/services/api-gateway/grpc-clients"
	"ride-sharing/shared/contracts"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if reqBody.UserId == "" {
		http.Error(w, "user id is required", http.StatusBadRequest)
		return
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		http.Error(w, "failed to marshal request body", http.StatusInternalServerError)
		return
	}

	reader := bytes.NewReader(jsonBody)

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		http.Error(w, "failed to create trip service client", http.StatusInternalServerError)
		return
	}
	defer tripService.Close()

	res, err := http.Post("http://trip-service:8083/preview", "application/json", reader)
	if err != nil {
		http.Error(w, "failed to send request to trip service", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	var responseBody any
	if err := json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		http.Error(w, "failed to decode response body", http.StatusInternalServerError)
		return
	}

	responses := contracts.APIResponse{Data: responseBody}
	fmt.Println("test from handTrip Preview")
	writeJSON(w, http.StatusOK, responses)
}
