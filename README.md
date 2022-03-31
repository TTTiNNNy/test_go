# Go Coding Challenge - J

## Testing (or getting) Knowledge

- Go-lang
- gRPC / gRPC-Streaming
- gRPC Metadata
- HTTP(REST)
- Go testing
- Docker
- Console and CLI tools

## Tasks

> Read through all the tasks and notes before solving.

### F | 12%pp

:heavy_check_mark: Fork this repo\
:heavy_check_mark: Clone -> Now you can push your changes\
:heavy_check_mark: Install [buf](https://docs.buf.build/introduction) and generate code from [proto](pkg/proto/challenge.proto)\
:heavy_check_mark: Make gRPC Server base: `main.go` in `cmd/server`, server methods in `pkg/server`

### D | 3%pp

:heavy_check_mark: Read environment vars: [task description](#environment)\
:heavy_check_mark: Make Metadata reader method\
:heavy_check_mark: Make link shortener method

### C | 3%pp

:heavy_check_mark: Make link metadata test(s)\
:heavy_check_mark: Make Timer Streaming method\
:heavy_check_mark: Make it running in Docker

### B | 5%pp

:heavy_check_mark: Make link shortener test(s)\
:heavy_check_mark: Make `cmd/client` using [cobra](https://github.com/spf13/cobra)

### A | 5%pp

- [ ] Make Timer Streaming test\
:heavy_check_mark: Gituhub - Publish, test and CI

## Important Notes to this tasks

1. Along with `cmd/server` - `cmd/client` has to be provided. For showing client capabilites write `go run main.go client -h`. For example, timer will start by using command `go run main.go client startTimer`. For additional option, use  `go run main.go client startTimer -h`
2. All test locate at test folder.

## Environment

There are no uasge for `BITLY_OAUTH_LOGIN` so I use only `BITLY_OAUTH_TOKEN` that using at `MakeShortLink` impl at `pkg/server/server.go`.
`main.go`

