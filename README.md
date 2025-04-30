# WebSocket Server with Channel Support

A high-performance WebSocket server implementation in Go that enables real-time communication through dynamic channels.
Built with Gorilla WebSocket and Mux libraries, this server provides a robust foundation for building scalable real-time
applications.

## Features

- **Dynamic Channel Support**: Create and join channels on-demand for isolated message broadcasting
- **Concurrent Client Handling**: Efficiently manages multiple concurrent connections using Go's goroutines
- **Custom Hostname Mapping**: Maps local IP to custom hostname for easier development and testing
- **Cross-Platform**: Full support for Windows, Linux, and macOS systems
- **Zero-Configuration**: Works out of the box with sensible defaults
- **Thread-Safe**: Uses proper synchronization for concurrent operations

---

## Getting Started

### Prerequisites

- **Go**: Version 1.24 or later
- **Required Libraries**: