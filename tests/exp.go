package tests

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

// Upgrader for WebSocket connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// Client represents a WebSocket connection
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Hub manages multiple WebSocket connections (per channel)
type Hub struct {
	clients map[string]map[*Client]bool
	lock    sync.RWMutex
}

var hub = Hub{
	clients: make(map[string]map[*Client]bool),
}

// Handle WebSocket connections
func handleConnections(w http.ResponseWriter, r *http.Request) {
	channel := mux.Vars(r)["channel"]

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	client := &Client{
		conn: conn,
		send: make(chan []byte),
	}

	// Add client to the hub
	hub.lock.Lock()
	if hub.clients[channel] == nil {
		hub.clients[channel] = make(map[*Client]bool)
	}
	hub.clients[channel][client] = true
	hub.lock.Unlock()

	// Handle sending and receiving
	go writeMessage(client)
	readMessage(client, channel)

	// Remove the client when done
	hub.lock.Lock()
	delete(hub.clients[channel], client)
	if len(hub.clients[channel]) == 0 {
		delete(hub.clients, channel)
	}
	hub.lock.Unlock()
	close(client.send)
}

// Read messages from the WebSocket connection
func readMessage(client *Client, channel string) {
	defer client.conn.Close()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		// Broadcast to other clients on the same channel
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

// Write messages to the WebSocket connection
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

// Get the machine's local IP address
func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// Skip interfaces that are down or loopback
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

			// Return the first valid IPv4 address
			if ip != nil && ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}

	return "", errors.New("no active network interfaces found")
}

// Get the path to the system hosts file based on OS
func getHostsFilePath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}

// Check if a hostname exists in the hosts file
func hostnameExistsInHosts(hostname string) (bool, error) {
	hostsFile := getHostsFilePath()

	file, err := os.Open(hostsFile)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), hostname) {
			return true, nil
		}
	}

	return false, nil
}

// Add a hostname-to-IP mapping to the hosts file
func addHostnameToHosts(hostname, ip string) error {
	hostsFile := getHostsFilePath()

	// Open the hosts file in append mode
	file, err := os.OpenFile(hostsFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the mapping in the correct format
	entry := fmt.Sprintf("%s %s\n", ip, hostname)
	_, err = file.WriteString(entry)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	port := flag.String("port", "8080", "Port to start the server on")
	hostname := flag.String("hostname", "socketflow", "Hostname to map to the local IP")
	flag.Parse()

	localIP, err := getLocalIP()
	if err != nil {
		log.Fatalf("Failed to determine local IP: %v", err)
	}

	// Try to check and add hostname to hosts file
	exists, err := hostnameExistsInHosts(*hostname)
	if err != nil {
		log.Printf("Error checking hosts file: %v", err)
	}
	if !exists {
		fmt.Printf("Hostname '%s' not found in hosts file, attempting to add it...\n", *hostname)

		err := addHostnameToHosts(*hostname, localIP)
		if err != nil {
			log.Printf("Failed to add hostname to hosts file: %v", err)
			fmt.Printf("Could not add hostname '%s' to hosts file. You may need to run the program with elevated privileges (e.g., 'sudo').\n", *hostname)
			fmt.Printf("Falling back to using the local IP address directly: %s\n", localIP)
			*hostname = localIP // Use the local IP directly if hostname can't be added
		} else {
			fmt.Printf("Added hostname '%s' mapped to '%s'.\n", *hostname, localIP)
		}
	} else {
		fmt.Printf("Hostname '%s' already exists in hosts file.\n", *hostname)
	}

	// Start the WebSocket server
	addr := fmt.Sprintf("%s:%s", localIP, *port)
	r := mux.NewRouter()
	r.HandleFunc("/{channel}", handleConnections)

	fmt.Printf("WebSocket server running at:\n")
	fmt.Printf("    Hostname: ws://%s:%s/{channel}\n", *hostname, *port)
	fmt.Printf("    Local IP: ws://%s:%s/{channel}\n", localIP, *port)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
