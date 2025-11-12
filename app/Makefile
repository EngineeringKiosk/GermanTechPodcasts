.DEFAULT_GOAL := help

.PHONY: help
help: ## Outputs the help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Compiles the application
	go build -race -o GermanTechPodcasts

.PHONY: staticcheck
staticcheck: ## Runs static code analyzer staticcheck
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./...

.PHONY: vet
vet: ## Runs go vet
	go vet ./...

.PHONY: tagstats
tagstats: ## get tag stats
	./GermanTechPodcasts tagStats --yml-directory ../podcasts/ 

.PHONY: run
run: build ## Compiles and starts the application
	./GermanTechPodcasts

.PHONY: test
test: ## Runs all unit tests
	go test -v -race ./...

.PHONY: extract-emails
extract-emails: build ## Extracts contact emails from RSS feeds and saves to JSON file
	./GermanTechPodcasts extractEmails --json-directory ../generated --output ../emails_with_authors.json
