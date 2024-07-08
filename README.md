# HTTP Client for Date-Time Service

This Go package provides functionalities to interact with a date-time service over HTTP. It includes methods to fetch the current date and time from a specified server endpoint, handle different response formats (JSON and plain text), and retry failed operations using an exponential backoff strategy.
## Installation

To install the project use:

```bash
go get github.com/codescalersinternships/datetime-client-amryassir
```

## Usage

Load configuration from environment variables or use default values
``` go
config := pkg.LoadConfig()
```

Initialize a new client instance
``` go
client := pkg.NewClient(config)
```
Get the current date and time from the server
``` go
dateTime, err := client.GetDateTime()
```

## Testing
Run the tests using Go's testing package.
```
go test ./...
```
## Contribution
Feel free to open issues or submit pull requests with improvements.

