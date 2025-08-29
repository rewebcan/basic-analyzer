# Home24 Web Page Analyzer

A web application built in Go that analyzes web pages and provides comprehensive insights about their structure, content, and accessibility.

## 🎯 Objective

Build a web application that performs analysis of web pages/URLs. The application provides a user-friendly interface where users can input a URL and receive detailed analysis results including HTML version, page title, heading structure, link analysis, and login form detection.

## ✨ Features

### Core Functionality
- **Web Form Interface**: Clean, simple form for URL input with submission button
- **URL Validation**: Comprehensive validation for HTTP/HTTPS URLs with proper error messages
- **Real-time Analysis**: Instant analysis results displayed after form submission

### Analysis Results
- **HTML Version Detection**: Identifies the HTML document version
- **Page Title Extraction**: Displays the page title from `<title>` tags
- **Heading Analysis**: Counts and categorizes headings by level (H1-H6)
- **Link Analysis**:
  - Internal vs external link identification
  - Link accessibility testing
  - Broken link detection and reporting
- **Login Form Detection**: Identifies pages containing password input fields

### Error Handling
- **HTTP Status Code Reporting**: Detailed error messages with status codes
- **Network Error Handling**: Graceful handling of unreachable URLs
- **Validation Errors**: Clear feedback for invalid URLs or malformed requests

## 🏗️ Architecture

### Project Structure
```
├── cmd/web/                 # Main application entry point
│   ├── main.go             # HTTP server setup and routing
│   └── tests/              # End-to-end tests
│       └── handler_e2e_test.go
├── internal/
│   ├── crawler/            # Crawling logic and HTTP handlers
│   │   ├── crawler.go      # Core crawling functionality
│   │   ├── handler.go      # HTTP request handlers
│   │   └── crawler_test.go # Unit tests
│   ├── fetcher/            # HTTP fetching and HTML parsing
│   │   ├── fetcher.go      # Fetcher implementation
│   │   ├── fetch_client.go # HTTP client utilities
│   │   └── fetcher_test.go # Unit tests
│   └── util/               # Utility functions
│       ├── config.go       # Configuration management
│       ├── url.go          # URL validation and normalization
│       └── types.go        # Common types
├── views/                  # HTML templates
│   └── index.html         # Main web interface
└── Makefile               # Build and test automation
```

### Design Decisions

#### 1. Modular Architecture
- **Separation of Concerns**: Clear separation between fetching, crawling, and presentation layers
- **Interface-based Design**: `Fetcher` and `Crawler` interfaces allow for easy testing and extension
- **Dependency Injection**: Constructor injection for better testability

#### 2. Concurrent Link Checking
- **Worker Pool Pattern**: Uses goroutines with semaphore for controlled concurrency
- **Configurable Limits**: Concurrency limits configurable via environment variables
- **Timeout Protection**: HTTP client timeouts prevent hanging requests

#### 3. Template Flexibility
- **Configurable Templates**: Support for custom template paths for testing
- **Backward Compatibility**: Maintains existing API while adding new features

#### 4. Error Handling Strategy
- **Structured Errors**: Custom error types with context
- **HTTP Status Mapping**: Proper HTTP status codes for different error scenarios
- **User-Friendly Messages**: Clear, actionable error messages for users

## 🚀 Installation & Setup

### Prerequisites
- Go 1.24.2 or later
- Git

### Quick Start
```bash
# Clone the repository
git clone <repository-url>
cd url-fetcher-home24

# Install dependencies
go mod download

# Run the application
make run
# or
go run cmd/web/

# The application will be available at http://localhost:8080/
```

### Building the Application

#### Build for Current Platform
```bash
# Build binary for your current platform
make build

# The binary will be created at ./build/url-fetcher
```

#### Cross-Platform Builds
```bash
# Build for Linux (amd64)
make build-linux

# Build for macOS (amd64)
make build-darwin

# Build for Windows (amd64)
make build-windows

# Build for all platforms at once
make build-all
```

#### Running the Built Binary
```bash
# After building, run the binary directly
./build/url-fetcher

# Or with custom configuration
CRAWLER_TIMEOUT=30s ./build/url-fetcher
```

#### Clean Build Artifacts
```bash
# Remove all build artifacts and binaries
make clean
```

### Configuration
The application supports configuration via environment variables:

```bash
# Crawler timeout (default: 10s)
export CRAWLER_TIMEOUT=30s

# Body size limit (default: 10MB)
export CRAWLER_BODY_SIZE_LIMIT=20971520

# Concurrency limit (default: 10)
export CRAWLER_CONCURRENCY_LIMIT=20
```

## 📖 Usage

1. **Access the Application**: Open http://localhost:8080/ in your web browser
2. **Enter URL**: Type a valid HTTP/HTTPS URL in the input field
3. **Submit**: Click the submit button to analyze the page
4. **View Results**: Review the comprehensive analysis results

### Example Usage
```
Input URL: https://example.com

Results:
- HTML Version: HTML 5
- Page Title: Example Domain
- Headings: H1 (1), H2 (0), H3 (0), H4 (0), H5 (0), H6 (0)
- Links: 2 internal, 1 external, 0 broken
- Login Form: No
```

## 🧪 Testing

### Unit Tests
Run unit tests for all packages:
```bash
make test
```

### End-to-End Tests
Run comprehensive e2e tests:
```bash
make e2e-test
```

### Test Coverage
- **Unit Tests**: Cover individual components and functions
- **E2E Tests**: Test complete user workflows with realistic scenarios
- **Error Scenarios**: Test validation, network errors, and edge cases

## 🔧 Development

### Code Quality
- **Linting**: Uses golangci-lint for code quality checks
- **Formatting**: Follows Go standard formatting
- **Testing**: Comprehensive test coverage with race detection

### Available Make Targets
```bash
make help          # Show all available targets
make run           # Run the application
make test          # Run unit tests
make e2e-test      # Run end-to-end tests
make lint          # Run linters
make lint-fix      # Auto-fix linting issues
make build         # Build for current platform
make build-linux   # Build for Linux (amd64)
make build-darwin  # Build for macOS (amd64)
make build-windows # Build for Windows (amd64)
make build-all     # Build for all platforms
make clean         # Clean build artifacts
```

## Possible Improvements
- **Robots.txt Compliance**: Check and respect robots.txt files before crawling to follow website rules
- **Rate Limiting**: Implement intelligent rate limiting to avoid overwhelming target servers
- **User-Agent Identification**: Proper user-agent strings identifying the crawler
- **Crawl Delay Respect**: Honor crawl-delay directives from robots.txt
- **Smart Concurrency**: Adaptive concurrency based on server response times to prevent DDoS-like behavior
- **Better Portability**: Using Docker to containarize the application for better portability.

## 📄 License

This project is developed as part of a technical assessment for Home24.

## 🤝 Contributing

This is a demonstration project. For contributions or questions, please refer to the project documentation.

---