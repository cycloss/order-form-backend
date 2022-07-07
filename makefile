
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
	(cd playground && go run .)

BASE = compose-files/docker-compose.base.yml

BASE_FLAGS = -f $(BASE) --project-directory .

DEV = compose-files/docker-compose.dev.yml

DEV_FLAGS = -f $(BASE) -f $(DEV) --project-directory .

up: ## start all dev services (without building)
	docker-compose $(DEV_FLAGS) up -d --remove-orphans

down: ## stop all dev services	
	docker-compose $(DEV_FLAGS) down --remove-orphans

build: ## build all dev services
	docker-compose $(DEV_FLAGS) build

TEST = compose-files/docker-compose.test.yml

TEST_FLAGS = -f $(BASE) -f $(TEST) --project-directory .

test: ## run all tests
	docker-compose $(TEST_FLAGS) up	--build --remove-orphans order-api-test

db-shell: ## start mysql client in db's uinvest database
	docker-compose $(BASE_FLAGS) exec db sh -c 'mysql -u$${MYSQL_USER} -p$${MYSQL_PASSWORD} $${MYSQL_DATABASE}'	

db-populate: ## populate the database
	export TEST_COMMAND="go test -v ./test -run \"^TestDbPopulate$$\"" && \
	docker-compose $(TEST_FLAGS) up --build --remove-orphans order-api-test

db-clear: ## clear the database
	export TEST_COMMAND="go test -v ./test -run \"^TestDbClear$$\"" && \
	docker-compose $(TEST_FLAGS) up --build --remove-orphans order-api-test