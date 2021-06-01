FROM golang:alpine as builder
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
RUN go mod download

#Copy code into the container
COPY . .

#Build the app
RUN go build -o main .

#Move to dist directory
WORKDIR /dist

#Copy binary from build
RUN cp /build/main .

FROM alpine
RUN apk add build-base
COPY --from=builder /dist/main .

#run
ENTRYPOINT ["./main"]