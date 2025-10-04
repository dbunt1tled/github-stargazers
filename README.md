# GitHub Stargazers Tracker

[![Go Report Card](https://goreportcard.com/badge/github.com/dbunt1tled/github-stargazers)](https://goreportcard.com/report/github.com/dbunt1tled/github-stargazers)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/dbunt1tled/github-stargazers.svg)](https://pkg.go.dev/github.com/dbunt1tled/github-stargazers)
[![GitHub release](https://img.shields.io/github/release/dbunt1tled/github-stargazers.svg)](https://github.com/dbunt1tled/github-stargazers/releases/)

A Go application that tracks stargazers for your GitHub repositories over time. It helps you monitor the growth and changes in your repository's stargazers.

## Features

- Track stargazers for all your GitHub repositories
- Store historical stargazer data in SQLite database
- Identify new and lost stargazers between runs
- Simple configuration with environment variables
- Lightweight and fast execution

## Prerequisites

- Go 1.25 or higher
- GitHub Personal Access Token with `public_repo` scope
- Git

## Installation

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/dbunt1tled/github-stargazers.git
   cd github-stargazers
   ```

2. Build the application:
   ```bash
   go build -o github-stargazers ./cmd/main.go
   ```

### Using Go Install

```bash
go install github.com/dbunt1tled/github-stargazers/cmd/github-stargazers@latest
```

## Configuration

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit the `.env` file with your GitHub credentials:
   ```
   GITHUB_USERNAME="your-github-username"
   GITHUB_TOKEN="your-github-token"
   DATABASE_PATH="./data.db"
   ```

## Usage

Run the application:

```bash
./github-stargazers
```

The application will:
1. Fetch all repositories for the specified GitHub user
2. Get current stargazers for each repository
3. Store the data in the SQLite database
4. Show new and lost stargazers compared to the previous run

## Project Structure

```
.
├── cmd/
│   └── main.go               # Main application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── db/                   # Database operations
│   └── github/               # GitHub API client
├── .env.example             # Example environment variables
├── go.mod                   # Go module definition
└── go.sum                   # Go module checksums
```

## Building and Running

### Build

```bash
go build -o github-stargazers ./cmd/main.go
```

### Run

```bash
./github-stargazers
```

### Run with custom config

```bash
GITHUB_USERNAME=your-username GITHUB_TOKEN=your-token ./github-stargazers
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [go-github](https://github.com/google/go-github) for GitHub API interactions
- Data stored in [SQLite](https://sqlite.org/)
