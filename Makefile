.PHONY: all
all: build_test cli web cdn api

.PHONY: build_test
build_test:
	cd cmd/cli && go build po.go && rm po
	cd cmd/api && go build main.go && rm main
	cd cmd/cdn && go build main.go && rm main
	cd examples/simple && go build main.go && rm main

.PHONY: web
web:
	cd ../podops.dev && gridsome build
	rm -rf cmd/cdn/public
	cp -R ../podops.dev/dist cmd/cdn/public

.PHONY: api
api:
	cd cmd/api && gcloud app deploy . --quiet

.PHONY: cdn
cdn:
	cd cmd/cdn && gcloud app deploy . --quiet

.PHONY: cli
cli:
	cd cmd/cli && go mod verify && go mod tidy && go install po.go