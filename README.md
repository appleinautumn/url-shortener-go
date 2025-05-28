# URL Shortener

This is a simple URL shortener service implemented in Go. It provides an HTTP API to shorten long URLs into short IDs and to resolve short IDs back to the original URLs. The service uses SQLite as its storage backend.

## Features

- Shorten long URLs to short IDs via a POST request.
- Resolve short IDs to the original URLs via a GET request.
- Stores URL mappings in a SQLite database.

## Getting Started

### Prerequisites

- Go 1.18 or later.
- SQLite3.

### Installation

1. Clone the repository:

   ```bash
   git clone git@github.com:appleinautumn/url-shortener-go.git
   cd url-shortener-go
   ```

2. Build the project:

   ```bash
   go build -o url-shortener cmd/main.go
   ```

3. Run the service:

   ```bash
   ./url-shortener
   ```

   The server listens on port `8080`.

### Running the Program

Make sure you have a `.env` file in the project root with the following content:

```
DB_FILE=dev.db
```

Then run the program executable:

```bash
./url-shortener
```

The server will start listening on port `8080`.

## API Endpoints

### Shorten URL

- **URL:** `/shorten`
- **Method:** `POST`
- **Request Body:**

  ```json
  {
    "long": "https://example.com/very/long/url"
  }
  ```

- **Response:**

  ```json
  {
    "short": "abc123",
    "long": "https://example.com/very/long/url"
  }
  ```

### Resolve Short URL

- **URL:** `/{shortID}`
- **Method:** `GET`
- **Response:**

  ```json
  {
    "short": "abc123",
    "long": "https://example.com/very/long/url"
  }
  ```

- Returns HTTP 404 if the short ID is not found.

## Project Structure

- `main.go`: The main application entry point, HTTP server, and routing.
- `storage/`: Package handling SQLite database initialization and URL storage/retrieval.
- `urls.db`: SQLite database file storing URL mappings.

## Notes

- The short ID is generated as a random hexadecimal string.
- The SQLite driver is imported and used within the `storage` package.
- Make sure to replace the module import path in `main.go` if your module name differs.

## License

This project is licensed under the MIT License.
