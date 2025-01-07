package sdk

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type CLIParams struct {
	Args     map[string]any
	ArgsList []string
	Flags    map[string]any
}

func ReceiveCLIParams() (*CLIParams, error) {
	ln, err := net.Listen("tcp", ":0") // Ephemeral port
	if err != nil {
		log.Fatalf("child: failed to listen: %v\n", err)
	}
	defer ln.Close()

	fmt.Printf("EPHEMERAL_ADDR=%s\n", ln.Addr().String())
	_ = os.Stdout.Sync()

	log.Printf("child: listening on %s\n", ln.Addr().String())

	// Accept a single connection (for demo)
	conn, err := ln.Accept()
	if err != nil {
		return nil, fmt.Errorf("child: accept error: %v", err)
	}
	defer conn.Close()

	// Read a line from parent, respond
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		// log.Printf("child: received from parent: %s", line)
		// fmt.Fprintf(conn, "child: got your message: %s\n", line)

		// Add small delay to ensure message is sent before closing
		time.Sleep(100 * time.Millisecond)

		// Parse the JSON message
		var params CLIParams
		if err := json.Unmarshal([]byte(scanner.Text()), &params); err != nil {
			return nil, fmt.Errorf("failed to unmarshal params: %v", err)
		}

		log.Println("child: successfully received and parsed params")
		return &params, nil
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("child: scanner error: %v", err)
	}

	return nil, fmt.Errorf("child: no line received")

}
