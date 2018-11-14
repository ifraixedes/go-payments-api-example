
define HELP_MSG
Execute one of the following targets:

endef

export HELP_MSG


.PHONY: help
help: ## Show this help
	@echo "$$HELP_MSG"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/:.*##/:##/' | column -t -s '##'

.PHONY: go-tools-install
go-tools-install: .gti-metalinter ## Install Go tools

.PHONY: lint
lint: ## Lint the code
	@gometalinter --vendor --enable-all --line-length=120 --warn-unmatched-nolint --exclude=vendor --exclude=mock_ --deadline=5m ./...

.PHONY: .go-tools-install-ci
.go-tools-install-ci: .gti-metalinter

.PHONY: .gti-metalinter
.gti-metalinter:
	@go get -u github.com/alecthomas/gometalinter
# gometalinter --install is deprecated but for now it's the best solution:
# https://github.com/alecthomas/gometalinter/issues/418
	@gometalinter --install --update --force --debug
