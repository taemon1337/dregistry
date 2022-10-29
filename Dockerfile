## Build Image
FROM golang:1.19-alpine AS build

WORKDIR /go/src/app

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

### Run Image
FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app /

EXPOSE 7946

ENTRYPOINT ["/app"]
