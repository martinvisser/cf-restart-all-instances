GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOCOVER=$(GOCMD) tool cover
GOCOVERCOBERTURA=gocover-cobertura
GOJUNITREPORT=go-junit-report
GOLINT=$$($(GOCMD) env GOPATH)/bin/staticcheck
GOARCH=amd64
PLATFORMS=linux darwin windows
MKDIR=mkdir -p
OUTPUT_DIRECTORY=bin/
BINARY_NAME=$(OUTPUT_DIRECTORY)cf-restart-all-instances
BUILD_DIRECTORY=build/
LDFLAGS=-ldflags="-X 'main.artifactVersion=${ARTIFACT_VERSION}' -s -w"
COVERAGE_FILE=$(BUILD_DIRECTORY)coverage.cov
COVERAGE_REPORT=$(BUILD_DIRECTORY)coverage.html
COVERAGE_REPORT_COBERTURA=$(BUILD_DIRECTORY)coverage.xml
TEST_REPORT=$(BUILD_DIRECTORY)report.xml
COVERAGE_OPTS=-coverprofile $(COVERAGE_FILE) 2>&1 | $(GOJUNITREPORT) > $(TEST_REPORT)
COVERAGE_EXPORT_OPTS=-html=$(COVERAGE_FILE) -o $(COVERAGE_REPORT)
CF_COMMAND=cf install-plugin

export SHELLOPTS:=pipefail

.PHONY: build test coverage lint clean $(addprefix build_,$(PLATFORMS)) $(addprefix install_,$(PLATFORMS))

all: lint test build install

build: $(addprefix build_,$(PLATFORMS))

$(addprefix build_,$(PLATFORMS)): build_%:
	GOOS=$(patsubst build_%,%,$@) GOARCH=$(GOARCH) $(GOBUILD) $(LDFLAGS) -v -o $(BINARY_NAME)-$(patsubst build_%,%,$@)
test:
	$(GOTEST) -v
coverage:
	$(MKDIR) $(BUILD_DIRECTORY)
	$(GOTEST) -v $(COVERAGE_OPTS)
	$(GOCOVER) $(COVERAGE_EXPORT_OPTS)
	$(GOCOVERCOBERTURA) < $(COVERAGE_FILE) > $(COVERAGE_REPORT_COBERTURA)
lint:
	$(GOLINT)
clean:
	$(GOCLEAN)
	rm -rf $(OUTPUT_DIRECTORY) $(BUILD_DIRECTORY)
install: build install_darwin
$(addprefix install_,$(PLATFORMS)): install_%:
	$(CF_COMMAND) -f $(BINARY_NAME)-$(patsubst install_%,%,$@)
