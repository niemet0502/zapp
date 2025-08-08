package connection

import "sync"

type UserConnected struct {
	ID   string
	Type string
	Chan chan string
}

var (
	Drivers   = make(map[string]UserConnected)
	DriversMu sync.Mutex
)
