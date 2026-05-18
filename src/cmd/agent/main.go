package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"ckrof-agent/internal/docker"
	"ckrof-agent/internal/models"
	"ckrof-agent/internal/queue"
	"ckrof-agent/internal/sysinfo"
	"ckrof-agent/internal/transport"

	"github.com/docker/docker/client"
)

func main() {
	// 1. Extract the RabbitMQ URL string
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("Fatal: RABBITMQ_URL environment variable is required and cannot be empty.")
	}

	// 2. Extract the Queue Name (with a safe fallback)
	queueName := os.Getenv("CKROF_QUEUE_NAME")
	if queueName == "" {
		queueName = "telemetry_queue"
	}
	authToken := os.Getenv("CKROF_SECRET_TOKEN")
	if authToken == "" {
		authToken = "ckrof-default-secret-token"
	}

	agentPort := os.Getenv("CKROF_AGENT_PORT")
	if agentPort == "" {
		agentPort = "9090" // Default port if not supplied
	}

	// 2. Parse telemetry sleep interval dynamically (with fallback to 30s)
	scanIntervalSeconds := 30
	if intervalStr := os.Getenv("CKROF_SCAN_INTERVAL"); intervalStr != "" {
		if parsed, err := strconv.Atoi(intervalStr); err == nil {
			scanIntervalSeconds = parsed
		}
	}
	scanInterval := time.Duration(scanIntervalSeconds) * time.Second

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Fatal error connecting to Docker: %v", err)
	}
	defer cli.Close()

	// 3. Start Control Server passing the environment-driven port
	go transport.StartControlServer(cli, authToken, agentPort)

	log.Printf("Ckrof Agent initialized. Telemetry interval set to %v.", scanInterval)

	for {
		// 1. Catch BOTH values from your sysinfo package
		fallbackHostname, fallbackIp := sysinfo.GetHostInfo()

		// 2. Set Hostname
		hostname := os.Getenv("CKROF_NODE_NAME")
		if hostname == "" {
			hostname = fallbackHostname
		}

		// 3. NEW: Set Real IP
		ipAddress := os.Getenv("CKROF_NODE_IP")
		if ipAddress == "" {
			ipAddress = fallbackIp
		}

		containers, err := docker.ScanAll(ctx)
		if err != nil {
			log.Printf("Error scanning Docker: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		payload := models.Payload{
			Host: models.HostInfo{
				Hostname: hostname,
				IP:       ip,
			},
			Timestamp:  time.Now().Format(time.RFC3339),
			Containers: containers,
		}

		err = queue.Publish(ctx, rabbitURL, payload, queueName)
		if err != nil {
			log.Printf("Failed to publish telemetry to RabbitMQ: %v", err)
		} else {
			log.Printf("Telemetry payload for %d container(s) pushed to queue.", len(containers))
		}

		// Use the flexible time interval
		time.Sleep(scanInterval)
	}
}