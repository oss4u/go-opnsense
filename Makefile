init:
	go mod tidy
#	go mod verify
#	go mod vendor

test:
	go test ./...

test_ci:
	export CI=true go test ./...
#	export OPNSENSE_ADDRESS="https://127.0.0.1:10443/" && \
#	export OPNSENSE_KEY="R+lLPklKa2QbfPtcpNeWwI9uNaDcd8ZRFJDUDpgH3uKvdyFn9HpOMqmsevTk5RDpk4FUjZFwgu2JHgQ5" && \
#	export OPNSENSE_SECRET="WE27qrbOZxTopTFZdsPdEt0rg8Uhqb6tmS44EuAkLsNk3oMY1GKvB4Zjp9S7oqEJeKM03QoH0QxFp76z" && \


start_vagrant:
	vagrant up

stop_vagrant:
	vagrant stop

rm_vagrant:
	vagrant destroy -f

build: init
	go build -o dist/opn-macos-amd64 opn.go
	env GOSS=linux go build -o dist/opn-linux-amd64 opn.go
	env GOSS=windows go build -o dist/opn-windows-amd64 opn.go
release:
	go build -o dist/opn-macos-amd64-${VERSION} opn.go
	env GOSS=linux go build -o dist/opn-linux-amd64-${VERSION} opn.go