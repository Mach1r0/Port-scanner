FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/scanner ./cmd/scanner

FROM alpine:3.23

RUN adduser -D -u 10001 scanner

COPY --from=build /out/scanner /usr/local/bin/scanner

USER scanner

ENTRYPOINT ["/usr/local/bin/scanner"]
