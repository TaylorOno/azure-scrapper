.PHONY: sonar-scanner
## sonar: *REQUIRES LOCAL SONARQUBE and SONAR-SCANNER run uploads to local sonar instance
sonar: sonar-scanner
	sonar-scanner -Dsonar.projectKey=trello-tribbles -Dsonar.exclusions=**/*_test.go,**/test_data/*,**/mocks/**,**/main.go -Dsonar.host.url=http://localhost:9000 -Dsonar.source=. -Dsonar.go.coverage.reportPaths=**/coverage.out

.PHONY: test
## test: runs all tests
test:
	go test ./... -coverprofile=./coverage.out

.PHONY: vet
## vet: runs go vet
vet:
	go vet ./...

.PHONY: fmt
## fmt: runs go fmt
fmt:
	go fmt ./...

.PHONY: pre-release
## pre-release: runs all tests and go tools
pre-release: fmt vet test

.PHONE: build
## Builds the binary for release
build: pre-release
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags' -a -o scrapper
	upx --brute scrapper

.PHONE: deploy
## Builds the binary for release
deploy: build


.PHONY: help
## help: Prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ":" | sed -e 's/^/ /'