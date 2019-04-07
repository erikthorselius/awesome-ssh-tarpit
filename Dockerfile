FROM golang:1.11-stretch as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o awesome-ssh-tarpit -mod vendor .
FROM scratch
COPY --from=builder /build/awesome-ssh-tarpit /app/
WORKDIR /app
CMD ["./awesome-ssh-tarpit"]