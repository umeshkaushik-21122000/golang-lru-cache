# LRU Cache with GET / SET / DELETE APIs in Go

This project implements a Least Recently Used (LRU) cache with GET, SET, and DELETE APIs in Go. It  includes a backend server that exposes these APIs as endpoints.

## Features

- LRU Cache implementation with adjustable expiration times for keys.
- Concurrent safe operations using Go's built-in concurrency mechanisms.
- Backend server built with Go and HTTP handlers for GET, SET, and DELETE operations.
- WebSocket implementation to dynamically reflect current key-value pairs and their expiration times in the UI.


### Backend Setup (Go)

1. Clone the repository:

   ```bash
   git clone <repository_url>
   cd lru-cache-golang
   ```

2. Install dependencies:
   ```bash
    go mod tidy
    ```
3. Run the backend Server:
   ```bash
    go run main.go
   ```
