DAGGER ?= dagger

init:
	$(DAGGER) run -- make _init

_init:
	go mod tidy
#	go mod verify
#	go mod vendor

test:
	$(DAGGER) run -- make _test

_test:
	go test ./...

test_ci:
	$(DAGGER) run -- make _test_ci

_test_ci:
	CI=true go test ./...
#	export OPNSENSE_ADDRESS="https://127.0.0.1:10443/" && \
#	export OPNSENSE_KEY="R+lLPklKa2QbfPtcpNeWwI9uNaDcd8ZRFJDUDpgH3uKvdyFn9HpOMqmsevTk5RDpk4FUjZFwgu2JHgQ5" && \
#	export OPNSENSE_SECRET="WE27qrbOZxTopTFZdsPdEt0rg8Uhqb6tmS44EuAkLsNk3oMY1GKvB4Zjp9S7oqEJeKM03QoH0QxFp76z" && \


start_vagrant:
	vagrant up

stop_vagrant:
	vagrant stop

rm_vagrant:
	vagrant destroy -f

build:
	$(DAGGER) run -- make _build

_build:
	go build -v ./...

ci:
	$(DAGGER) run -- make _ci

_ci:
	go mod tidy
	go build -v ./...
	go test -v ./...

release:
	$(DAGGER) run -- make _release

_release:
	go build -v ./...