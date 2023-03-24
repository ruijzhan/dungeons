FROM golang:1.20 as builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o server cmd/server/main.go

FROM alpine:3.17
WORKDIR /
COPY --from=builder /workspace/server .
USER 65532:65532

ENTRYPOINT [ "/server" ]
