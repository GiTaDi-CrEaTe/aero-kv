package wal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/GiTaDi-CrEaTe/aero-kv/core"
)

// Logger handles the append-only Write-Ahead Log for durability.
type Logger struct {
	mu   sync.Mutex
	file *os.File
}

// NewLogger opens or creates the WAL file.
func NewLogger(filename string) (*Logger, error) {
	// os.O_APPEND means we only ever add to the bottom of the file
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: file}, nil
}

// AppendSet records a SET mutation to the disk.
func (l *Logger) AppendSet(key string, value []byte, ttl time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var expiresAt int64
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl).UnixNano()
	}

	// Format: SET,key,value,expiresAt
	entry := fmt.Sprintf("SET,%s,%s,%d\n", key, string(value), expiresAt)
	l.file.WriteString(entry)
	l.file.Sync() // CRITICAL: Force OS to flush buffer to physical disk
}

// AppendDel records a DEL mutation to the disk.
func (l *Logger) AppendDel(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := fmt.Sprintf("DEL,%s\n", key)
	l.file.WriteString(entry)
	l.file.Sync()
}

// Restore reads the WAL on boot and rebuilds the memory engine state.
func Restore(filename string, store *core.Store) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No WAL found. Starting fresh.")
			return
		}
		fmt.Println("CRITICAL: Failed to read WAL:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	loaded := 0

	// Replay the history line by line
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) == 0 {
			continue
		}

		cmd := parts[0]
		if cmd == "SET" && len(parts) >= 4 {
			key := parts[1]
			value := []byte(parts[2])
			expiresAt, _ := strconv.ParseInt(parts[3], 10, 64)

			var ttl time.Duration
			if expiresAt > 0 {
				now := time.Now().UnixNano()
				if expiresAt <= now {
					continue // This key already died in the past, skip it
				}
				ttl = time.Duration(expiresAt - now)
			}
			store.Set(key, value, ttl)
			loaded++
		} else if cmd == "DEL" && len(parts) >= 2 {
			key := parts[1]
			store.Delete(key)
			loaded--
		}
	}

	fmt.Printf("WAL Restore Complete: Rebuilt state with %d historical operations.\n", loaded)
}