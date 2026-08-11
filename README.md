# eBECAS V2 Client

A Go client for interacting with the eBECAS V2 API.

## Requirements

- Go 1.26.5 or later
- An eBECAS V2 API access token

## Installation

```bash
go get github.com/skills-australia-institute/ebecas-v2-client
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	ebecasv2client "github.com/skills-australia-institute/ebecas-v2-client"
)

func main() {
	client := ebecasv2client.NewClient(
		ebecasv2client.Config{
			BaseURL:     "https://college.ap1.ebecas.app/api/v2",
			AccessToken: "YOUR_API_ACCESS_TOKEN",
			PageSize:    20,
		},
	)

	ctx := context.Background()

	user, statusCode, err := client.Me(ctx)
	if err != nil {
		log.Fatalf("request failed with status %d: %v", statusCode, err)
	}

	fmt.Println(statusCode)
	fmt.Println(user)
}
```

## Configuration

Create a client using `NewClient` and provide a `Config`:

```go
client, err := ebecasv2client.NewClient(
	ebecasv2client.Config{
		BaseURL:     "https://college.ap1.ebecas.app/api/v2",
		AccessToken: "YOUR_API_ACCESS_TOKEN",
		PageSize:    20,
	},
)
if err != nil {
	log.Fatal(err)
}
```

### Configuration Options

| Option        | Required | Default           | Description                                                                                                                 |
| ------------- | -------- | ----------------- | --------------------------------------------------------------------------------------------------------------------------- |
| `BaseURL`     | Yes      | —                 | Base URL of the eBECAS V2 API. It must include the `/api/v2` path.                                                          |
| `AccessToken` | Yes      | —                 | eBECAS V2 API access token used for Bearer authentication.                                                                  |
| `HTTPClient`  | No       | 15-second timeout | Optional custom `*http.Client`. If `nil`, a client with a 15-second timeout is used.                                        |
| `PageSize`    | No       | `10`              | Number of records requested per API page. Must be between `1` and `100`. A value of `0` uses the default page size of `10`. |

### Custom HTTP Client

You can provide a custom `http.Client` when you need a different timeout or HTTP configuration:

```go
httpClient := &http.Client{
	Timeout: 30 * time.Second,
}

client, err := ebecasv2client.NewClient(
	ebecasv2client.Config{
		BaseURL:     "https://college.ap1.ebecas.app/api/v2",
		AccessToken: "YOUR_API_ACCESS_TOKEN",
		HTTPClient:  httpClient,
		PageSize:    20,
	},
)
if err != nil {
	log.Fatal(err)
}
```

## Authentication

The client uses Bearer token authentication:

```http
Authorization: Bearer YOUR_API_ACCESS_TOKEN
```

Do not commit API tokens or credentials to source control. Use environment variables or a secure secrets-management solution.

For example:

```go
accessToken := os.Getenv("EBECAS_ACCESS_TOKEN")
```

## Error Handling

Client methods return the API response, HTTP status code, and error:

```go
response, statusCode, err := client.SomeMethod(ctx)
```

Always check the returned error before using the response:

```go
student, statusCode, err := client.CreateStudent(ctx, input)
if err != nil {
	log.Printf("request failed with status %d: %v", statusCode, err)
	return
}

fmt.Println(student)
```

## Context

API methods accept a `context.Context`, allowing requests to be cancelled or given a timeout.

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

user, statusCode, err := client.Me(ctx)
```

## License

This project is licensed under the [MIT License](https://github.com/skills-australia-institute/ebecas-v2-client/blob/main/LICENSE).
