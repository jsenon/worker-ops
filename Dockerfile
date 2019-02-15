FROM golang:1.11.5 as build-env

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /
CMD ["/app"]