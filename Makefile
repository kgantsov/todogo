run_go_dev:
	swag init -g ./cmd/server/main.go -o ./docs; cd cmd/server/; go run main.go

build_go:
	swag init -g ./cmd/server/main.go -o ./docs; cd cmd/server/; go build


build:
	go test ./... -cover
	swag init -g ./cmd/server/main.go -o ./docs; cd cmd/server/; go build

build_linux:
	go test ./... -cover
	swag init -g ./cmd/server/main.go -o ./docs; cd cmd/server/; GOOS=linux GOARCH=amd64 GOOS=linux GOARCH=amd64 go build

test:
	go test ./... -cover

load_test:
	k6 run --summary-trend-stats="med,p(95),p(99.9)" load_testing/script.js

benchmark:
	cd load_testing; plow ${UPTIME_HOST}/API/v1/services -c 100 -n 10000 -T 'application/json' -m GET -H "Authorization: Bearer ${UPTIME_TOKEN}"
