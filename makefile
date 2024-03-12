# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# Install dependencies

gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

brew:
	brew update
	brew list encore || brew install encoredev/tap/encore

# ==============================================================================
# Run Project

up:
	encore run -v

GENERATE_ID = $(shell docker ps | grep encoredotdev | cut -c1-12)
SET_ID = $(eval MY_ID=$(GENERATE_ID))

down:
	$(SET_ID)
	docker stop $(MY_ID)
	docker rm $(MY_ID) -v

# ==============================================================================
# Access Project

token:
	curl -il \
	--user "admin@example.com:gophers" http://localhost:3000/v1/users/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

users:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/v1/users?page=1&rows=2"

pgcli:
	pgcli $(shell encore db conn-uri url | sed -e 's/localhost/127.0.0.1/g')

curl:
	curl -il "http://127.0.0.1:4000/test?limit=2&offset=2"
