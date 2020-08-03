FROM golang:1.14-alpine AS build

WORKDIR /app/
COPY cmd/proxy/proxy.go /app/cmd/proxy/
COPY internal/utils/utils.go /app/internal/utils/
COPY go.* /app/
RUN CGO_ENABLED=0 go build -o gcrproxy ./cmd/proxy

FROM scratch
COPY --from=build /app/gcrproxy /bin/gcrproxy
ENTRYPOINT ["/bin/gcrproxy"]
