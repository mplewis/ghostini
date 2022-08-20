FROM golang:1.19 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

ENV CGO_ENABLED=0
COPY . .
RUN go build -o /build/app .

########################################

FROM scratch
WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/app .

ENTRYPOINT ["/app/app"]
