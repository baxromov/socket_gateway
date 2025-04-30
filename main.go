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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients map[string]map[*Client]bool
	lock    sync.RWMutex
}

var hub = Hub{
	clients: make(map[string]map[*Client]bool),
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	channel := mux.Vars(r)["channel"]

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

	hub.lock.Lock()
	if hub.clients[channel] == nil {
		hub.clients[channel] = make(map[*Client]bool)
	}
	hub.clients[channel][client] = true
	hub.lock.Unlock()

	go writeMessage(client)
	readMessage(client, channel)

	hub.lock.Lock()
	delete(hub.clients[channel], client)
	if len(hub.clients[channel]) == 0 {
		delete(hub.clients, channel)
	}
	hub.lock.Unlock()
	close(client.send)
}

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

func getHostsFilePath() string {
	if runtime.GOOS == "windows" {
		return `C:\Windows\System32\drivers\etc\hosts`
	}
	return "/etc/hosts"
}

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

func addHostnameToHosts(hostname, ip string) error {
	hostsFile := getHostsFilePath()

	file, err := os.OpenFile(hostsFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

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
			*hostname = localIP
		} else {
			fmt.Printf("Added hostname '%s' mapped to '%s'.\n", *hostname, localIP)
		}
	} else {
		fmt.Printf("Hostname '%s' already exists in hosts file.\n", *hostname)
	}

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
