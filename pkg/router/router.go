package router

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Tubular-Bytes/statesman/pkg/backend"
	"github.com/Tubular-Bytes/statesman/pkg/model"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(map[string]any{"status": "healthy"}); err != nil {
		slog.Error("failed to encode health response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func HandleLock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	defer r.Body.Close()

	lockData, err := getLockFromRequest(r)
	if err != nil {
		slog.Error("failed to get lock data from request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	store := backend.Get()

	if err := store.Lock(lockData); err != nil {
		http.Error(w, "resource is locked", http.StatusConflict)

		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(lockData); err != nil {
		slog.Error("failed to encode lock data", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func HandleUnlock(w http.ResponseWriter, r *http.Request) {
	lockData, err := getLockFromRequest(r)
	if err != nil {
		slog.Error("failed to get lock data from request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	store := backend.Get()
	if err := store.Unlock(lockData.LockID); err != nil {
		slog.Error("failed to unlock", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	slog.Debug("lock dropped", "lockID", lockData.LockID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(lockData); err != nil {
		slog.Error("failed to encode lock data", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func HandlePostState(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("ID")
	slog.Debug("state post", "id", id)

	var state model.State

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&state); err != nil {
		slog.Error("failed to decode request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	store := backend.Get()
	if err := store.PutState(id, &state); err != nil {
		slog.Error("failed to put state", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(state); err != nil {
		slog.Error("failed to encode state", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func HandleGetState(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("ID")

	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(map[string]any{"version": 1}); err != nil {
			slog.Error("failed to encode state", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		return
	}

	slog.Debug("state get", "id", id)

	store := backend.Get()

	state, err := store.GetState(id)
	if err != nil {
		slog.Error("failed to get state", "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(state); err != nil {
		slog.Error("failed to encode state", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
}

func getLockFromRequest(r *http.Request) (*model.LockData, error) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	var lockData model.LockData

	if err := decoder.Decode(&lockData); err != nil {
		slog.Error("failed to decode request body", "error", err)

		return nil, err
	}

	return &lockData, nil
}
