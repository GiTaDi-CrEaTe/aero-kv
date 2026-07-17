package main

import (
	"fmt"
	"time"

	"github.com/GiTaDi-CrEaTe/aero-kv/core"
	"github.com/GiTaDi-CrEaTe/aero-kv/server"
	"github.com/GiTaDi-CrEaTe/aero-kv/wal"
)

func main() {
	// 1. Boot up the concurrent memory engine
	db := core.NewStore()

	// 2. Replay historical data from the WAL (if it exists)
	wal.Restore("aero.wal", db)

	// 3. Initialize the Logger to record new mutations
	logger, err := wal.NewLogger("aero.wal")
	if err != nil {
		fmt.Println("CRITICAL: Failed to initialize WAL")
		panic(err)
	}

	// 4. Start the Garbage Collector
	db.StartSweeper(5 * time.Second)

	// 5. Open the TCP gates
	server.Start("9000", db, logger)
}