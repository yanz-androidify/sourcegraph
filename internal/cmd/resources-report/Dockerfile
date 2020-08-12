FROM golang:1.15-alpine@sha256:59eae48746048266891b7839f7bb9ac54a05cec6170f17ed9f4fd60b331b644b AS builder

WORKDIR /go/src/resources-report

COPY go.* ./
RUN go mod download

COPY *.go ./
RUN go build -o /bin/resources-report

FROM alpine:3.11@sha256:cb8a924afdf0229ef7515d9e5b3024e23b3eb03ddbba287f4a19c6ac90b8d221

# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/resources-report /usr/local/bin/

ENTRYPOINT ["resources-report"]
