.PHONY: help run build db-migrate db-down db-fix

TEST_FLAGS ?=
GOOGLE_APPLICATION_CREDENTIALS := $$PWD/config/google-credentials.json
UPPER_DB_LOG=DEBUG
CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

all:
	@echo "*******************************"
	@echo "**   functions build tool    **"
	@echo "*******************************"
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  run                   - run in dev mode"
	@echo ""
	@echo "  deploy                - deploy to production"
	@echo ""
	@echo "  tail-prod             - tail logs form production"
	@echo ""
	@echo ""

print-%: ; @echo $*=$($*)


##
## Tools
##
tools:
	@mkdir -p ./bin
	@go mod tidy

conf:
	[ -f config/api.conf ] || cp config/api.develop.conf config/api.conf

.PHONY: run
run: conf
	@(FUNCTION_TARGET=HelloWorld LOCAL_ONLY=true go run cmd/functions/main.go)

.PHONY: tail-prod
tail-prod:
	@echo "do nothing"

.PHONY: deploy
deploy:
	@read -r -p "Are you sure you want to deploy to production? y/n: " confirm; \
	case $${confirm} in \
		y)     \
			git push prod master \
			;; \
		n) \
			echo done; \
			;; \
		*) \
			echo please enter y/n; \
			;; \
	esac;
