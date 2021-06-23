
# Set these according to your environment.
# TODO: Can we get these from an env var?

# Your IPA Host url.  NOTE: No trailing slash
IPA_HOST=https://ipa.dev.iad2.dc.paas.redhat.com

# Your ipa_session cookie.
# This will change regularly.  See `developing.md` for how to get this value.
COOKIE=MagBearerToken=WuAc5m2kymwQ6mh1Vbu7xEBv23YjjNCwyxGcxpJWnqt265j7cJIdIMcWt1pF9yOJl3S1ytiOk4ADZmvOUTXF6dzI5IpxrjXJLX2AB%2bRSPAMgGLgAZOQAhk5C%2b4waZl7hJLIPxngXdSrl9u%2bv69Bio2qx2Y40F%2fEYjOfP0M9UToQBv8UeV7Tr018n0mOYynh9bKs1V6N2ANCHHokQRxgUtA%3d%3d


SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec



##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: dump-schema
dump-schema: ## Download the schema json file from the IPA server
	$(shell ./scripts/dump-schema.sh --host ${IPA_HOST} --cookie ${COOKIE} --out ./data/schema.json)
	jq . ./data/schema.json > ./data/schema-new.json
	rm ./data/schema.json
	mv ./data/schema-new.json ./data/schema.json


.PHONY: dump-errors
dump-errors: ## Download the errors json file from GitHub
	python3 ./scripts/dump-errors.py ./data/errors.json
	jq . ./data/errors.json > ./data/errors-new.json
	rm ./data/errors.json
	mv ./data/errors-new.json ./data/errors.json


.PHONY: gen
gen: ## Regenerate the code for the client.
	cd gen && go run ./...
	go fmt ./freeipa/generated.go

