FROM golang:1.23-alpine as go-builder
ENV GIN_MODE=release
WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags 'libpq' -a -installsuffix cgo -ldflags="-w -s" -o tezos-delegation-service .
RUN ls -l

FROM scratch as go-runtime-container
ENV GIN_MODE=release

COPY --from=go-builder /app/tezos-delegation-service /tezos-delegation-service
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 3000 3001
ENTRYPOINT ["/tezos-delegation-service"]