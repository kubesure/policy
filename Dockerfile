FROM golang:1.12.4-alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 go install

FROM scratch
WORKDIR /opt
COPY --from=builder /go/bin/app .
ENTRYPOINT ["/opt/app"]
