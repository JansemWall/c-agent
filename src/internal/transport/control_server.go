package transport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ControlRequest struct {
	ContainerName string `json:"container_name"`
	Action        string `json:"action"`
}

// Accept the port parameter dynamically
func StartControlServer(cli *client.Client, token string, port string) {
	http.HandleFunc("/container/control", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Ckrof-Token") != token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ControlRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		var err error

		switch req.Action {
		case "start":
			err = cli.ContainerStart(ctx, req.ContainerName, container.StartOptions{})
		case "stop":
			err = cli.ContainerStop(ctx, req.ContainerName, container.StopOptions{})
		case "restart":
			err = cli.ContainerRestart(ctx, req.ContainerName, container.StopOptions{})
		default:
			http.Error(w, "Invalid action", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	})

	log.Printf("Control Server HTTP listening on port :%s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Fatal error on HTTP control server: %v", err)
	}
}