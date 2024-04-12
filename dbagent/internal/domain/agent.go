package domain

import "time"

type Agent struct {
	Name            string    `json:"agent_unique_name"`
	CreatedAt       time.Time `json:"created_at"`
	LastHeartbeatAt time.Time `json:"last_heartbeat_at"`
}
