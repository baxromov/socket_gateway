
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
    - `ws://localhost:5000/mychannel`
    - `ws://socketflow:5000/chatroom`

---

### WebSocket Client Examples

#### **Using JavaScript**

Clients can use the WebSocket API in their browsers to connect to the server.
```javascript
// Establish a connection to the WebSocket server
const channel = "mychannel"; // Replace with your channel name
const ws = new WebSocket(`ws://localhost:5000/${channel}`);

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
   -H "Host: localhost:5000" -H "Origin: http://localhost:5000" \
   ws://localhost:5000/mychannel
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
wscat -c ws://localhost:5000/mychannel
```
Once connected, you can type messages in the terminal, and they will be broadcast to other clients in the same channel!

---

## Future Improvements

- Add support for secure WebSocket (WSS) protocols with built-in TLS handling.
- Implement authentication for channels to restrict access.
- Provide Docker support for easy containerization and deployment.
- Add advanced monitoring and logging.

---

This documentation outlines all necessary details, including setup, usage examples, and deployment instructions for your WebSocket Server. Let me know if you need additional sections or further clarification!
```
baxromov.shahzodbek@gmail.com
```
