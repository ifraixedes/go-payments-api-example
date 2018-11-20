OS := "osx"
OS_ALT := $(shell echo $(OS) | sed 's/osx/darwin/')
MODD_VER := "0.7"
NODE_VER := "11.2.0"


define HELP_MSG
Requirements:
	Some of the targets requires to have docker installed.

Development workflow for working on the openapi spec and see the swagger
documentation on each update:
	* Execute the 'swagger-ui-up' target
	* In one terminal instance, execute the 'openapi-merge-watch' target
	* In another terminal instance, execute the 'swagger-ui-autorefresh' target

	Then you can changes the partial openapi spec files and they will be
	automatically merged and rendered by the swagger UI. Check each command help
	message for knowing the http address to load.

Execute one of the following targets:

endef

export HELP_MSG


.PHONY: help
help: ## Show this help
	@echo "$$HELP_MSG"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/:.*##/:##/' | column -t -s '##'

.PHONY: install-tools
install-tools: ## Install tools to work with the API specification
	@mkdir -p .bin
	@echo
	@echo "Installing modd"
	@echo
	@curl -L -o .bin/modd.tgz \
		https://github.com/cortesi/modd/releases/download/v$(MODD_VER)/modd-$(MODD_VER)-$(OS)64.tgz
	@tar Uxzf .bin/modd.tgz -C .bin
	@mv -f .bin/modd-$(MODD_VER)-$(OS)64/modd .bin/
	@rm -fr .bin/modd.tgz .bin/modd-$(MODD_VER)-$(OS)64
	@echo
	@echo "Installing Node"
	@echo
	@curl -L -o .bin/node.tgz https://nodejs.org/dist/v$(NODE_VER)/node-v$(NODE_VER)-$(OS_ALT)-x64.tar.gz
	@tar Uxzf .bin/node.tgz -C .bin
	@mv .bin/node-v$(NODE_VER)-$(OS_ALT)-x64 .bin/node
	@rm -f .bin/node.tgz
	@make install-tools-update-node-deps

.PHONY: install-tools-update-node-deps
install-tools-update-node-deps: ## Update the NodeJS packages used as tooling
	@PATH=.bin/node/bin npm install

.PHONY: openapi-merge
openapi-merge: ## Merge the partial openapi json files into openapi-merged.json file
	@PATH=.bin/node/bin:$$PATH ./node_modules/.bin/swagger-merger -i openapi.json -o openapi-merged.json

.PHONY: openapi-merge-watch
openapi-merge-watch: ## Like openapi-merge but running a process which watches changes in the partil files and update the merged file on each change.
	@.bin/modd

.PHONY: swagger-ui-up
swagger-ui-up: .swagger-ui-down ## Runs swagger-ui server (http://localhost:8081)
	@docker run -d -p 8081:8080 -e SWAGGER_JSON=/tmp/openapi-merged.json -v `pwd`:/tmp --name swagger-ui-payments swaggerapi/swagger-ui > /dev/null

.PHONY: swagger-ui-autorefresh
swagger-ui-autorefresh: ## Runs a proxy server (http://localhost:8080) for swagger-ui server which autoreload the browser when openapi-merged.json file changes
	@PATH=.bin/node/bin:$$PATH ./node_modules/.bin/browser-sync  'http://localhost:8081' 'openapi-merged.json' -w --no-open --port 8080 --no-ui

.PHONY: clean
clean: .swagger-ui-down ## Clean up the environment (docker containers, etc.)

.PHONY: .swagger-ui-down
.swagger-ui-down:
	@docker rm -f swagger-ui-payments 2>&1>/dev/null || true
