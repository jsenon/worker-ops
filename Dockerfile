FROM golang:latest as build

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


FROM gcr.io/distroless/base
COPY --from=build /app/worker-ops /
ENV PORT=8080
EXPOSE 8080
CMD ["/worker-ops"]
