
# Getube

Getube is a Go project that uses the `youtube.Client` to stream videos directly to the client. The project uses the `go.mod` for managing dependencies.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

You need to have Go installed on your machine. You can download it from the official [Go website](https://golang.org/dl/).

### Installing

Clone the repository to your local machine:

```bash
git clone https://github.com/nguyenvanduocit/getube.git
cd getube
go mod download
```

## Running the Application

To run the application, use the following command:

```bash
go run main.go
```

The server will be spinning on port `8080` by default. You can change the port by setting the `PORT` environment variable.

Open your browser and go to `http://localhost:8080/:videoID` to see the application in action.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* This project is for educational purposes only.
