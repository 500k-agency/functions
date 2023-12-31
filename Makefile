.PHONY: help

TEST_FLAGS ?=
EXAMPLE_CONFIG := $$PWD/config/function.example.conf
LOCAL_CONFIG := $$PWD/config/function.conf
GOOGLE_APPLICATION_CREDENTIALS := $$PWD/config/google-credentials.json
UPPER_DB_LOG=DEBUG
CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

all:
	@echo "*******************************"
	@echo "**   project build tool      **"
	@echo "*******************************"
	@echo "make <cmd>"
	@echo ""
	@echo "commands:"
	@echo "  run                   - run functions in dev mode"
	@echo "  toolkit               - run toolkit to initialize project setup"
	@echo ""
	@echo "  deploy                - deploy to production"
	@echo ""
	@echo "  tail-prod             - tail logs form production"
	@echo ""
	@echo ""

print-%: ; @echo $*=$($*)

bins := gcloud go
$(bins):
	@which $@ > /dev/null

.PHONY: install-gcloud
install-gcloud:
	@if ! which gcloud > /dev/null 2>&1; then\
		read -r -p "Have you installed gcloud cli tool?? y/n: " confirm; \
		case $${confirm} in \
			y)     \
				;; \
			n) \
				echo "Please install via https://cloud.google.com/sdk/docs/install"; \
				;; \
			*) \
				echo please enter y/n; \
				;; \
		esac; \
    fi

##
## Tools
##
.PHONY: tools
tools: install-gcloud
	@mkdir -p ./bin
	@go mod tidy

.PHONY: toolkit
toolkit:
	@($(GO) run cmd/toolkit/main.go -config=${LOCAL_CONFIG})

.PHONY: conf
conf:
	@[ -f ${LOCAL_CONFIG} ] || cp ${EXAMPLE_CONFIG} ${LOCAL_CONFIG}

.PHONY: run
run: conf
	@(FUNCTION_TARGET=PurchaseHandler LOCAL_ONLY=true go run cmd/functions/main.go -config=${LOCAL_CONFIG})

.PHONY: tail-prod
tail-prod:
	@echo "do nothing"

.PHONY: deploy
deploy:
	@read -r -p "Are you sure you want to deploy function to production? y/n: " confirm; \
	case $${confirm} in \
		y)     \
		gcloud functions deploy DevPurchaseHandler \
			--trigger-http \
			--runtime=go121 \
			--entry-point=PurchaseHandler \
			--set-secrets '/etc/secrets/latest=CLOUDFN_SECRET_DEV:latest' \
			;; \
		n) \
			echo done; \
			;; \
		*) \
			echo please enter y/n; \
			;; \
	esac;
