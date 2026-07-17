package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/GiTaDi-CrEaTe/aero-kv/core"
	"github.com/GiTaDi-CrEaTe/aero-kv/wal"
)

// Start boots up the TCP server and now takes the WAL Logger
func Start(port string, store *core.Store, logger *wal.Logger) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("CRITICAL: Failed to bind to port", port)
		return
	}
	defer listener.Close()

	fmt.Println("🚀 Aero-KV is airborne. Listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, store, logger)
	}
}

func handleConnection(conn net.Conn, store *core.Store, logger *wal.Logger) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {
		case "PING":
			conn.Write([]byte("PONG\n"))

		case "SET":
			if len(parts) < 3 {
				conn.Write([]byte("ERROR: SET requires key and value\n"))
				continue
			}
			key := parts[1]
			value := []byte(parts[2])
			var ttl time.Duration
			if len(parts) == 4 {
				sec, err := strconv.Atoi(parts[3])
				if err == nil {
					ttl = time.Duration(sec) * time.Second
				}
			}
			
			// 1. Save to RAM
			store.Set(key, value, ttl)
			// 2. Safely log to physical disk
			logger.AppendSet(key, value, ttl)
			
			conn.Write([]byte("OK\n"))

		case "GET":
			if len(parts) < 2 {
				conn.Write([]byte("ERROR: GET requires a key\n"))
				continue
			}
			key := parts[1]
			val, exists := store.Get(key)
			if !exists {
				conn.Write([]byte("(nil)\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("%s\n", string(val))))
			}

		case "DEL":
			if len(parts) < 2 {
				conn.Write([]byte("ERROR: DEL requires a key\n"))
				continue
			}
			key := parts[1]
			
			// 1. Delete from RAM
			store.Delete(key)
			// 2. Log the deletion
			logger.AppendDel(key)
			
			conn.Write([]byte("OK\n"))

		default:
			conn.Write([]byte("ERROR: Unknown command\n"))
		}
	}
}