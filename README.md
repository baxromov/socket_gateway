```markdown
# WebSocket Server with Channel Support

A high-performance WebSocket server implementation in Go that enables real-time communication through dynamic channels.  
Utilizing the powerful Gorilla WebSocket and Mux libraries, the server is designed to support scalable, real-time applications efficiently.

---

## Features

- **Dynamic Channel Support**: Dynamically create isolated communication channels for broadcasting messages.
- **Concurrent Client Management**: Seamlessly manage multiple WebSocket connections using Go's goroutines.
- **Custom Hostname Mapping**: Automatically maps local IP to a custom hostname or uses sensible defaults for development and testing.
- **Cross-Platform**: Fully compatible with Windows, Linux, and macOS environments.
- **Zero-Configuration**: Works out of the box with sensible default configurations.
- **Thread-Safe**: Ensures proper synchronization for concurrent operations to avoid race conditions.

---

## Getting Started

### Prerequisites

- **Go**: Ensure Go is installed (version 1.24 or later).  
  To install Go, visit [Go Downloads](https://go.dev/dl/).
- **Required Libraries**: Install the necessary libraries listed below. Use the following commands to get the dependencies:
  ```bash
  go get github.com/gorilla/mux
  go get github.com/gorilla/websocket
  ```

---

### Installation

1. **Clone the Repository**:
   Clone the repository where this WebSocket project resides using the following command:
   ```bash
   git clone <repository_url>
   cd <repository_directory>
   ```

2. **Run the Server**:
   Execute the Go file to start the WebSocket server:
   ```bash
   go run main.go
   ```

---

### Configuration Options

You can customize the server using the available command-line flags:

- **`-port`**: Specify the port on which the server should run (default: `8080`).
- **`-hostname`**: Specify the custom hostname to map to the local IP (default: `socketflow`).

For example:
```
bash
go run main.go -port=3000 -hostname=mycustomhostname
```
The server would then listen at:
```

ws://mycustomhostname:3000/{channel}
```
---

### How It Works

1. **Dynamic Channels**:  
   The server handles WebSocket connections on dynamic, channel-based endpoints (e.g., `/{channel}`).  
   Channels allow client groups to communicate without interfering with each other.

2. **Concurrent Client Management**:  
   By using Go's goroutines, the server can handle many WebSocket connections simultaneously, providing efficient performance even under heavy load.

3. **Custom Hostname**:  
   The program automatically resolves the machine's local IP and attempts to map it to a hostname using the `hosts` file. If no mapping is found, the default hostname (`socketflow`) is used.

---

## Usage Examples

### Basic Example

- **WebSocket URL**:  
  Each channel on this WebSocket server can be accessed via the URL:
   ```
   ws://<hostname>:<port>/{channel}
   ```
  Replace `<hostname>`, `<port>`, and `{channel}` with your custom hostname, port, and channel name.

- Sample WebSocket endpoints:
    - `ws://localhost:8080/mychannel`
    - `ws://socketflow:8080/chatroom`

---

### WebSocket Client Examples

#### **Using JavaScript**

Clients can use the WebSocket API in their browsers to connect to the server.
```
javascript
// Establish a connection to the WebSocket server
const channel = "mychannel"; // Replace with your channel name
const ws = new WebSocket(`ws://localhost:8080/${channel}`);

// Listen for incoming messages
ws.onmessage = (event) => {
  console.log("Received message from server:", event.data);
};

// Send a message to the server
ws.onopen = () => {
  ws.send("Hello, WebSocket server!");
};

// Handle connection errors
ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};

// Handle connection closure
ws.onclose = () => {
  console.log("WebSocket connection closed.");
};
```
---

#### **Using `curl`**

You can interact with the WebSocket server using `curl`:

1. **Connect to the WebSocket Server**:
   ```bash
   curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" \
   -H "Host: localhost:8080" -H "Origin: http://localhost:8080" \
   ws://localhost:8080/mychannel
   ```

2. **Send Messages**:  
   While connected, you can manually send messages to the WebSocket server within your console or terminal.

---

#### **Using `wscat`**

`wscat` is a WebSocket client CLI tool. First, install it with npm:
```
bash
npm install -g wscat
```
Then, connect to your WebSocket server:
```
bash
wscat -c ws://localhost:8080/mychannel
```
Once connected, you can type messages in the terminal, and they will be broadcast to other clients in the same channel!

---

## Testing Locally

Here’s how you can test the WebSocket server:

1. Start the WebSocket server:
   ```bash
   go run main.go -port=8080
   ```

2. Open multiple browser tabs or terminals to connect clients using the WebSocket URL:
    - Browser JavaScript (DevTools Console):
      ```javascript
      const ws = new WebSocket("ws://localhost:8080/testchannel");
      ws.onmessage = (e) => console.log(e.data);
      ws.onopen = () => ws.send("Client connected!");
      ```
    - `wscat` or `curl` in a terminal.

3. Messages sent by one client will be broadcast to all other clients connected to the same channel.

---

## Deployment Notes

- **Reverse Proxy**: For production, deploy the server behind a reverse proxy like Nginx or Apache to manage SSL/TLS encryption.  
  Example Nginx config:
  ```nginx
  server {
      listen 443 ssl;
      server_name yourdomain.com;

      ssl_certificate /path/to/fullchain.pem;
      ssl_certificate_key /path/to/privkey.pem;

      location / {
          proxy_pass http://localhost:8080;
          proxy_http_version 1.1;
          proxy_set_header Upgrade $http_upgrade;
          proxy_set_header Connection "upgrade";
          proxy_set_header Host $host;
          proxy_cache_bypass $http_upgrade;
      }
  }
  ```

- **Port Binding**: Ensure the WebSocket server binds to a port that is accessible (e.g., `80` or `443` behind a proxy).

- **Environment Variables**: Consider using environment variables for port and hostname settings in production environments.

---

## Troubleshooting

1. **Port Already in Use**:  
   If the port is already in use, specify a different port with the `-port` flag:
   ```bash
   go run main.go -port=3000
   ```

2. **WebSocket Connection Fails**:
    - Ensure the server is running and accessible.
    - Check firewall settings to ensure the port is open.
    - If using a reverse proxy, verify its WebSocket configuration.

3. **Dynamic Hostname Not Working**:
    - Check your `/etc/hosts` file (or equivalent on Windows) to ensure the IP is properly mapped to a hostname.

---

## License

This project is open-source under the [MIT License](https://opensource.org/licenses/MIT). Feel free to use, modify, and distribute it for personal and commercial projects.

---

## Future Improvements

- Add support for secure WebSocket (WSS) protocols with built-in TLS handling.
- Implement authentication for channels to restrict access.
- Provide Docker support for easy containerization and deployment.
- Add advanced monitoring and logging.

---

This documentation outlines all necessary details, including setup, usage examples, and deployment instructions for your WebSocket Server. Let me know if you need additional sections or further clarification!
```
