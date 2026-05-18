package docker

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"ckrof-agent/internal/models"
)

func ScanAll(ctx context.Context) ([]models.ContainerInfo, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var results []models.ContainerInfo

	for _, c := range containers {
		var portMap []models.PortInfo
		for _, p := range c.Ports {
			portMap = append(portMap, models.PortInfo{
				Private: p.PrivatePort,
				Public:  p.PublicPort,
				Type:    p.Type,
			})
		}

		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
			if len(name) > 0 && name[0] == '/' {
				name = name[1:]
			}
		}

		var memUsage uint64
		stats, err := cli.ContainerStatsOneShot(ctx, c.ID)
		if err == nil {
			defer stats.Body.Close() 

			var raw struct {
				MemoryStats struct {
					Usage uint64 `json:"usage"`
				} `json:"memory_stats"`
			}

			// Decodifica o buffer de leitura diretamente
			if err := json.NewDecoder(stats.Body).Decode(&raw); err == nil {
				memUsage = raw.MemoryStats.Usage
			}
		}

		results = append(results, models.ContainerInfo{
			Name:        name,
			StackName:   c.Labels["com.docker.compose.project"],
			ServiceName: c.Labels["com.docker.compose.service"],
			Image:       c.Image,
			State:       c.State,
			Labels:      c.Labels,
			Ports:       portMap,
			MemUsage:    memUsage,
		})
	}

	return results, nil
}