package handlers

import (
	"encoding/json"
	"net/http"

	"casino_loyalty_reward_system/internal/types"

	"go.uber.org/zap"
)

const ResponseOK = "ok"

func WriteError(log *zap.SugaredLogger, w http.ResponseWriter, statusCode int, err error) {
	var errMessage string
	if err != nil {
		errMessage = err.Error()
	}
	WriteErrorMessage(log.WithOptions(zap.AddCallerSkip(1)), w, statusCode, errMessage)
}

func WriteErrorMessage(log *zap.SugaredLogger, w http.ResponseWriter, statusCode int, message string) {
	log.WithOptions(zap.AddCallerSkip(1)).Error(message)
	WriteJSON(log, w, statusCode, types.ErrorResponse{Message: message})
}

func WriteJSON(log *zap.SugaredLogger, w http.ResponseWriter, statusCode int, data any) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Errorf("failed to encode data: %s", err)
	}
}
