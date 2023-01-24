# build binary stage
FROM golang:1.19.5 AS builder

ARG REPO

LABEL stage=build-intermediate

WORKDIR ${REPO}

COPY . .
RUN mkdir -p /go/bin/
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/fls-interpreter ./cmd

# copy the server binary from builder stage; run the server binary
FROM alpine:latest AS runner

WORKDIR /bin
COPY --from=builder /go/bin/fls-interpreter .
ENTRYPOINT ["fls-interpreter"]
