package proxy

import (
	"log"
	"time"
)

type Watchdog struct {
	Timeout          int64
	ShutdownCallback func()
	lastUpdate       int64
}

func (watchdog *Watchdog) Reset() {
	watchdog.lastUpdate = time.Now().Unix()
}

func (watchdog *Watchdog) Start() {
	log.Println("Shutdown watchdog started.")

	for {
		status, err := GetServerStatus()
		if err == nil {
			someoneOnline := status.CurPlayers > 0

			if someoneOnline {
				watchdog.lastUpdate = time.Now().Unix()
			} else if time.Now().Unix()-watchdog.lastUpdate > watchdog.Timeout {
				log.Println("No one is online, shutting down...")
				watchdog.ShutdownCallback()
				watchdog.lastUpdate = time.Now().Unix()
			}

		}
		time.Sleep(time.Second * 2)
	}

}
