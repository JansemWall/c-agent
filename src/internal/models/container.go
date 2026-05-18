package models

type PortInfo struct {
	Private uint16 `json:"private"`
	Public  uint16 `json:"public"`
	Type    string `json:"type"`
}

type ContainerInfo struct {
	Name           string            `json:"name"`
	StackName      string            `json:"stack_name"`
	ServiceName    string            `json:"service_name"`
	Image          string            `json:"image"`
	State          string            `json:"state"`
	Ports          []PortInfo        `json:"ports"`
	Labels         map[string]string `json:"labels"`
	MemUsage       uint64            `json:"mem_usage"`
}

type HostInfo struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

type Payload struct {
	Host       HostInfo        `json:"host_info"`
	Timestamp  string          `json:"timestamp"`
	Containers []ContainerInfo `json:"containers"`
}