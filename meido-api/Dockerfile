ARG flask_url=https://meido-app.cf/flask/message_twitter_3_fast

FROM golang:latest as builder
WORKDIR /workdir
COPY . .
# Set Environment Variable
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV FLASK_URL=$flask_url

# Build
RUN go build -o app

# Runtime Container
FROM alpine
ARG flask_url
ENV FLASK_URL=$flask_url
RUN apk add --no-cache ca-certificates
COPY --from=builder /workdir/app /app
COPY *.txt /

ENTRYPOINT ["/app"]
