# FIRST STAGE
FROM golang:1.18 AS builder

# dependency management files
COPY share/go.* /share/
COPY order-api/go.* /app/

WORKDIR /app
# should only run mod download if any files in order-api .go files change
RUN go mod download

# copy shared libraries
COPY share/*.go /share/

COPY order-api/*.go ./

# compile, variables MUST be ARG and not ENV.
# CGO_ENABLED is mandatory as alpine doesn't have the dylib files required by the compiled binary
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/orderapi

# SECOND STAGE
FROM alpine:3.14.2
WORKDIR /app

COPY --from=builder /app/bin/orderapi /app

ENTRYPOINT [ "./orderapi" ]
