FROM node:22-alpine AS frontend

WORKDIR /app/web/frontend

COPY web/frontend/package.json web/frontend/package-lock.json ./
RUN npm ci

COPY web/frontend/ .
RUN npm run build

FROM golang:1.24-alpine AS builder

ARG VERSION=dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o velero-backup-reporter ./cmd/velero-backup-reporter/

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/velero-backup-reporter /usr/local/bin/velero-backup-reporter

USER 65534:65534

ENTRYPOINT ["velero-backup-reporter"]
