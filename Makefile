.PHONY: build frontend backend clean

build: frontend backend

frontend:
	cd web/frontend && npm ci && npm run build

backend:
	CGO_ENABLED=0 go build -o velero-backup-reporter ./cmd/velero-backup-reporter/

clean:
	rm -rf web/dist velero-backup-reporter
