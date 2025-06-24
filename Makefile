ifdef VER
VERSION=$(shell echo $(VER) | sed -e 's/^v//g' -e 's/\//_/g')
else
VERSION=$(shell grep -oP 'version\s*=\s*"\K[^"]+' flake.nix | head -n 1)
endif


ifeq ($(shell uname -m),x86_64)
  HOST_ARCH?=x86_64
  ARCH?=amd64
else ifeq ($(shell uname -m),arm64)
  HOST_ARCH?=aarch64
  ARCH?=arm64
else ifeq ($(shell uname -m),aarch64)
  HOST_ARCH?=aarch64
  ARCH?=aarch64
endif

ifeq ($(shell uname -o),Darwin)
  OS?=darwin
else
  OS?=linux
endif


.PHONY: help
help: ## Show this help.
	@IFS=$$'\n' ; \
	lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for line in $${lines[@]}; do \
		IFS=$$'#' ; \
		split=($$line) ; \
		command=`echo $${split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		info=`echo $${split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf "%-38s %s\n" $$command $$info ; \
	done


.PHONY: get-version
get-version:  ## Return version
	@sed -i '/^\s*version = "0\.0\.0-dev";/s//version = "${VERSION}";/' flake.nix
	@sed -i '/^\s*created = "1970-.*";/s//created = "${shell date --utc '+%Y-%m-%dT%H:%M:%SZ'}";/' flake.nix
	@echo $(VERSION)


.PHONY: check
check:   ## Run nix flake check
	./build/nix.sh flake check --print-build-logs


.PHONY: check-dry-run
check-dry-run:  ## Run nix flake check
	@nix build \
		--dry-run \
		--json \
		--print-build-logs \
		.\#checks.$(HOST_ARCH)-$(OS).go-checks | jq -r '.[].outputs.out'


.PHONY: dev-env-up
dev-env-up:  ## Starts development environment
	echo "Nothing to do"


.PHONY: dev-env-up-short
dev-env-up-short:  ## Starts development environment without ai service
	echo "Nothing to do"


.PHONY: dev-env-down
dev-env-down:  ## Stops development environment
	echo "Nothing to do"


.PHONY: build
build:  ## Build application and places the binary under ./result/bin
	nix build $(docker-build-options) \
		.\#mcp-nhost-$(ARCH)-$(OS) \
		--print-build-logs


.PHONY: build-dry-run
build-dry-run:  ## Run nix flake check
	@nix path-info \
		--derivation \
		.\#packages.$(ARCH)-$(OS).default


.PHONY: build-docker-image
build-docker-image:  ## Build docker image
	nix build $(docker-build-options) \
		.\#packages.$(HOST_ARCH)-linux.docker-image-$(ARCH) \
		--print-build-logs
	skopeo copy --insecure-policy \
		--override-arch $(ARCH) \
		dir:./result docker-daemon:nhost/mcp-nhost:$(VERSION)
