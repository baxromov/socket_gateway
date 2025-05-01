package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Upgrader configuration for WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// WebSocket Client structure
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Central hub for managing WebSocket clients
type Hub struct {
	clients map[string]map[*Client]bool
	lock    sync.RWMutex
}

// Create an instance of the Hub
var hub = Hub{
	clients: make(map[string]map[*Client]bool),
}

// Function to handle incoming WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Extract channel and subchannel from URL path parameters
	channel := mux.Vars(r)["channel"]
	subChannel := mux.Vars(r)["subchannel"] // Extract the additional subchannel path

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	log.Printf("New connection on channel: %s, subchannel: %s", channel, subChannel)

	client := &Client{
		conn: conn,
		send: make(chan []byte),
	}

	// Register client under the "channel/subChannel" namespace
	hub.lock.Lock()
	fullChannel := fmt.Sprintf("%s/%s", channel, subChannel) // Create a unique namespace
	if hub.clients[fullChannel] == nil {
		hub.clients[fullChannel] = make(map[*Client]bool)
	}
	hub.clients[fullChannel][client] = true
	hub.lock.Unlock()

	go writeMessage(client)
	readMessage(client, fullChannel)

	// When client disconnects, clean up
	hub.lock.Lock()
	delete(hub.clients[fullChannel], client)
	if len(hub.clients[fullChannel]) == 0 {
		delete(hub.clients, fullChannel)
	}
	hub.lock.Unlock()
	close(client.send)
}

// Function to read messages from a WebSocket connection
func readMessage(client *Client, channel string) {
	defer client.conn.Close()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		hub.lock.RLock()
		for subscriber := range hub.clients[channel] {
			select {
			case subscriber.send <- message:
			default:
				log.Println("Subscriber disconnected")
			}
		}
		hub.lock.RUnlock()
	}
}

// Function to write messages to a WebSocket connection
func writeMessage(client *Client) {
	defer client.conn.Close()

	for message := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

// Function to retrieve the local IP address of the machine
func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// Only consider interfaces that are up and not loopback
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Return the first IPv4 address found
			if ip != nil && ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}

	return "", errors.New("no active network interfaces found")
}

// Get the file path of the hosts file based on the platform
func getHostsFilePath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}

// Find the hostname in the /etc/hosts file for the provided IP
func getHostnameForIP(ip string) (string, error) {
	hostsFilePath := getHostsFilePath()

	file, err := os.Open(hostsFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip commented or empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[0] == ip {
			return parts[1], nil // Return the hostname found for the IP
		}
	}
	return "", errors.New("hostname not found for IP in hosts file")
}

// Main function
func main() {
	// Flags to specify port and optionally a hostname
	port := flag.String("port", "5000", "Port to start the server on")
	defaultHostname := flag.String("hostname", "socketflow", "Default hostname to map to the local IP")
	flag.Parse()

	// Get the local IP address
	localIP, err := getLocalIP()
	if err != nil {
		log.Fatalf("Failed to determine local IP: %v", err)
	}

	hostname, err := getHostnameForIP(localIP)
	if err != nil {
		log.Printf("Error finding hostname for local IP: %v", err)
		log.Printf("Falling back to the default hostname: %s\n", *defaultHostname)
		hostname = *defaultHostname
	} else {
		fmt.Printf("Found hostname '%s' for IP '%s'.\n", hostname, localIP)
	}

	addr := fmt.Sprintf("%s:%s", localIP, *port)
	r := mux.NewRouter()

	// Add new route that includes subchannel
	r.HandleFunc("/{channel}/.../{subchannel}", handleConnections)

	fmt.Printf("WebSocket server running at:\n")
	fmt.Printf("    Hostname: ws://%s:%s/{channel}/.../{subchannel}\n", hostname, *port)
	fmt.Printf("    Local IP: ws://%s:%s/{channel}/.../{subchannel}\n", localIP, *port)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
