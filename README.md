# Simple_gRPC-Go

This is a simple implementation of gRPC in Go. It includes basic CRUD operations on a `Product` entity.

## Requirements

- gRPC
- Protobuf

## Installation

1. Clone this repository:

```bash
git clone https://github.com/yourusername/Simple_gRPC-Go.git
```

2. Navigate into the project directory:

```
cd Simple_gRPC-Go
```

3. Download the dependencies:

```
go mod tidy
```

## Usage

1. Start the server:

```
go run cmd/main.go
```

The server should now be running and listening on port *50051*.

2. You can send requests to the server using a gRPC client. See protos/product.proto for the service and message definitions.

## Project Structure

- `cmd/main.go`: The main entry point for the server.
- `services/`: Contains the gRPC service implementations.
- `protos/`: Contains the Protobuf service and message definitions.
