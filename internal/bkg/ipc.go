package bkg

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

const (
	SocketFile = ".tyn/daemon.sock"
)

type Message struct {
	Command string          `json:"command"`
	Params  json.RawMessage `json:"params"`
}

type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error,omitempty"`
}

func SendCommand(cmd string, params interface{}) (*Response, error) {
	var paramsJSON []byte
	var err error
	if params != nil {
		paramsJSON, err = json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("error marshaling params: %w", err)
		}
	}

	msg := Message{
		Command: cmd,
		Params:  paramsJSON,
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("error marshaling message: %w", err)
	}

	sockPath, err := getSocketPath()
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("unix", sockPath)
	if err != nil {
		return nil, fmt.Errorf("error connecting to daemon: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(msgJSON)
	if err != nil {
		return nil, fmt.Errorf("error sending message to daemon: %w", err)
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading response from daemon: %w", err)
	}

	var resp Response
	err = json.Unmarshal(buf[:n], &resp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &resp, nil
}

func HandleConnections(handler func(Message) Response) error {
	sockPath, err := getSocketPath()
	if err != nil {
		return err
	}

	_, err = os.Stat(sockPath)
	if err == nil {
		err = os.Remove(sockPath)
		if err != nil {
			return fmt.Errorf("error removing existing socket: %w", err)
		}
	}

	sockDir := filepath.Dir(sockPath)
	err = os.MkdirAll(sockDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating socket directory: %w", err)
	}

	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		return fmt.Errorf("error creating socket listener: %w", err)
	}

	err = os.Chmod(sockPath, 0600)
	if err != nil {
		return fmt.Errorf("error setting socket permissions: %w", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
				continue
			}

			go handleConnection(conn, handler)
		}
	}()

	return nil
}

func handleConnection(conn net.Conn, handler func(Message) Response) {
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from connection: %v\n", err)
		return
	}

	var msg Message
	err = json.Unmarshal(buf[:n], &msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling message: %v\n", err)
		return
	}

	resp := handler(msg)

	respJSON, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		return
	}

	_, err = conn.Write(respJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing response: %v\n", err)
		return
	}
}

func getSocketPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	return filepath.Join(home, SocketFile), nil
}
