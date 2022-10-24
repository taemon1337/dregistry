## Build Image
FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /dregistry

### Run Image
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /dregistry /dregistry

EXPOSE 7373

USER nonroot:nonroot

ENTRYPOINT ["/dregistry"]
