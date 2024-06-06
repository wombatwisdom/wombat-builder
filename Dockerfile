FROM node:16.13.0 AS ui-build

WORKDIR /app
COPY internal/api/web .

RUN npm install
RUN npm run build

FROM golang:1.22.2 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /go/src/app

# Update dependencies: On unchanged dependencies, cached layer will be reused
COPY go.* /go/src/app
RUN go mod download

# Build
COPY . /go/src/app/
COPY --from=ui-build /app/dist internal/api/web/dist
RUN go build -o wombat-builder

# Pack
FROM debian:bookworm-slim AS package

LABEL maintainer="Daan Gerits <daan+wombat@shono.io>"

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/src/app/wombat-builder .

ENTRYPOINT ["/wombat-builder"]

CMD ["all", "--ui", "--port", "80"]

EXPOSE 80