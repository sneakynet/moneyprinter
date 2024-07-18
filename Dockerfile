FROM golang:1.22-alpine AS build
COPY . .
RUN apk add tini-static sqlite-dev build-base && CGO_ENABLED=1 go build -o /moneyd -ldflags '-extldflags "-static"' ./cmd/moneyd/main.go

FROM scratch
COPY --from=build /sbin/tini-static /tini
COPY --from=build /moneyd /moneyd
ENTRYPOINT ["/tini", "/moneyd"]


