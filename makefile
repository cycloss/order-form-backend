
# catch trying to run make if docker not running
# expands to nothing if not running
ifeq ($(docker info 2> /dev/null), "")
	$(error Docker must be running to use this makefile)
endif

.DEFAULT_GOAL := help
.PHONY: help playground devup devdown devrestart

help: ## Show this help
# grep grabs target lines using a regex from the MAKEFILE_LIST which is the name of this file
# double dollar is required to escape a dollar sign (end of line in regex)
# FS is the field separator which is ':', 'anything', then '## '
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) |\
	awk 'BEGIN {FS = ":.*?## "}; {printf "%-30s %s\n", $$1, $$2}'

certgen: ## generate self signed tls cert for nginx
	./nginx/certs/cert-gen.sh

playground: ## run the playground file
	go run playground/main.go

devup: ## start all services in the dev profiles
	docker-compose --profile dev up -d

devdown: ## stop all services in the dev profiles	
	docker-compose --profile dev down

devrestart: ## restart all services in the dev profiles	
	docker-compose --profile dev restart