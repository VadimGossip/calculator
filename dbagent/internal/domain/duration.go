package domain

import "time"

type OperationDuration struct {
	Name      string    `json:"operation_name"`
	Duration  uint16    `json:"operation_duration_ms"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
