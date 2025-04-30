# WebSocket Server with Channel Support

This project provides a simple WebSocket-based server that supports dynamic channels for real-time communication. The server is powered by Gorilla WebSocket and Mux libraries, delivering secure and scalable WebSocket connections for multiple clients.

## Features
- **Channel-based communication**: Users can join dynamic channels and broadcast messages within these channels.
- **Hostname setup**: Automatically maps a user-defined hostname to the local IP, making it easier to use custom hostnames for WebSocket testing.
- **Cross-platform support**: Works on both Windows and Unix-based systems.
- **Concurrent handling**: Supports multiple clients in different channels via Goroutines and synchronization mechanisms.

---

## Getting Started

### Prerequisites
- **Go**: Version 1.24 or later.
- **Gorilla WebSocket**: Ensure that the following libraries are installed before building this project:
```textmate
go get -u github.com/gorilla/mux
  go get -u github.com/gorilla/websocket
```

- **Administrator Rights**: If the program needs to modify the `/etc/hosts` file (Linux/macOS) or `C:\Windows\System32\drivers\etc\hosts` (Windows), it may require elevated privileges.

---

### Installation
1. Clone the repository or copy the source code into a new Go module:
```textmate
git clone <repository-url>
   cd <repository-folder>
```


2. Install the required dependencies:
```textmate
go mod tidy
```


3. Build the executable:
```textmate
go build -o websocket-server
```


---

### Usage

You can run the server using the following command:

```textmate
./websocket-server [flags]
```


#### Flags
- `--port`: The port on which the server starts. Default is `8080`.
- `--hostname`: The custom hostname the WebSocket server maps to the local IP. Default is `socketflow`.

#### Example
```textmate
./websocket-server --port 9090 --hostname myserver.local
```


---

### How It Works
- The server dynamically creates or removes "channels" based on client connections.
- Clients connect to specific channels via the WebSocket URL:  
  `ws://<hostname>:<port>/<channel>`.  
  For example:
    - `ws://socketflow:8080/chat` connects to the `chat` channel.
    - `ws://192.168.1.100:8080/updates` connects to the `updates` channel.
- Messages sent to a channel are broadcast to all clients connected to the same channel.

---

### Examples

1. **Start the server:**
```textmate
./websocket-server --port 8080 --hostname mycustomhost
```

Output:
```
WebSocket server running at:
       Hostname: ws://mycustomhost:8080/{channel}
       Local IP: ws://192.168.1.100:8080/{channel}
```


2. **Connect using a WebSocket client (e.g., `wscat`):**
```textmate
wscat -c ws://mycustomhost:8080/chat
```


3. **Broadcast messages in real-time:**
   If multiple clients connect to the `chat` channel, any message sent by one client will be broadcast to all other clients in that channel.

---

### Notes
- If the hostname is not in the system's `hosts` file, the program attempts to add it. If the process fails due to permission errors, the fallback is to use the local IP directly.
- For Windows users, ensure administrative permissions when modifying system files.

---

### Contribution
Feel free to fork this project and raise pull requests for:
- Adding features like authentication.
- Better error handling for WebSocket connections.
- Enhanced logging mechanisms.

For issues, open a GitHub issue or contact the repository owner.

---

### License
This project is licensed under the MIT License.