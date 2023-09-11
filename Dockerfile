FROM golang:1.21 as build

WORKDIR /build

# Layer so we don't have to download every time
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o main ./main.go

# Deploy
FROM gcr.io/distroless/base-debian12

COPY --from=build /build/main /main

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/main"]
