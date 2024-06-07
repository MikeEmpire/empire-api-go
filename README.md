# Empire API

## Production URL

<https://empire-api.afam.app>

Click [here](https://empire-api.afam.app) to view the app

## Overview

This is a scalable API built using the Gin web framework in Go. The project is structured to promote clean code, maintainability, and scalability.

## Project Structure

```bash
empire-api-go/
├── go.mod
├── go.sum
├── main.go
├── .env
├── pkg/
│ └── mail/
| |----helpers.go
| |----service.go
├── routes/
│ └── routes.go
└── config/
└─--- config.go
```

- **main.go**: The entry point of the application.
- **pkg/handlers**: Contains the handler functions for different routes.
- **routes/routes.go**: Sets up the routes for the application.
- **config/config.go**: Handles configuration settings.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (1.16+)
- [Gin](https://github.com/gin-gonic/gin)
- [GoDotEnv](https://github.com/joho/godotenv) (for configuration)

### Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/yourusername/your_project_name.git
   ```

2. Navigate to the project directory:

   ```sh
   cd your_project_name
   ```

3. Initialize the Go module:

   ```sh
   go mod init your_project_name
   ```

4. Install the dependencies:

   ```sh
   go get -u github.com/gin-gonic/gin
   go get github.com/joho/godotenv
   ```

### Configuration

Create a `.env` file in the root of the project to manage your environment variables:

### Steps to get credentials to use Gmail API

1. Make sure you have credentials saved in the root folder as `credentials.json`.

   Here is a link to the instructions to get the credentials: [Instructions](https://developers.google.com/gmail/api/quickstart/go).

2. Once you have the credentials saved in the root folder, run the `api/v1/auth` endpoint.

3. Get the authorization code and copy and paste it into the terminal and press enter. From there, the `token.json` file will be saved, and things will be good to go.

### Steps to getting Gmail Password

The [Instructions](https://support.google.com/accounts/answer/185833?hl=en&sjid=5919955758469844792-NC) to generate the app password and set it in the .env file

Please refer to the sample.env file for more information on how to configure your environment

## License

None