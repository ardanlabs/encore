# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# Deploy First Mentality

# ==============================================================================
# Go Installation
#
#	You need to have Go version 1.22 to run this code.
#
#	https://go.dev/dl/
#
#	If you are not allowed to update your Go frontend, you can install
#	and use a 1.22 frontend.
#
#	$ go install golang.org/dl/go1.22@latest
#	$ go1.22 download
#
#	This means you need to use `go1.22` instead of `go` for any command
#	using the Go frontend tooling from the makefile.

# ==============================================================================
# Brew Installation
#
#	Having brew installed will simplify the process of installing all the tooling.
#
#	Run this command to install brew on your machine. This works for Linux, Mac and Windows.
#	The script explains what it will do and then pauses before it does it.
#	$ /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
#
#	WINDOWS MACHINES
#	These are extra things you will most likely need to do after installing brew
#
# 	Run these three commands in your terminal to add Homebrew to your PATH:
# 	Replace <name> with your username.
#	$ echo '# Set PATH, MANPATH, etc., for Homebrew.' >> /home/<name>/.profile
#	$ echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> /home/<name>/.profile
#	$ eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
#
# 	Install Homebrew's dependencies:
#	$ sudo apt-get install build-essential
#
# 	Install GCC:
#	$ brew install gcc

# ==============================================================================
# Install Tooling and Dependencies
#
#	This project uses Docker and it is expected to be installed. Please provide
#	Docker at least 4 CPUs. To use Podman instead please alias Docker CLI to
#	Podman CLI or symlink the Docker socket to the Podman socket. More
#	information on migrating from Docker to Podman can be found at
#	https://podman-desktop.io/docs/migrating-from-docker.
#
#	Run these commands to install everything needed.
#	$ make dev-brew
#	$ make dev-gotooling

# ==============================================================================
# Running Test
#
#	Running the tests is a good way to verify you have installed most of the
#	dependencies properly.
#
#	$ make test
#

# ==============================================================================
# Running The Project
#
#	$ make up
#	$ make token
#	$ export TOKEN=<token>
#	$ make users
#
#	Use can use this command to shut everything down
#
#	$ make down

# ==============================================================================
# Deploy The Project
#
#	$ encore app init
#	Run the Cloud Dashboard
#		Github:  Connect the repo in the dashboard
#		Deploy:  Manually start a deploy
#       Secrets: make secrets

# ==============================================================================
# CLASS NOTES
#
# RSA Keys
# 	To generate a private/public key PEM file.
# 	$ openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# 	$ openssl rsa -pubout -in private.pem -out public.pem
#
# OPA Playground
# 	https://play.openpolicyagent.org/
# 	https://academy.styra.com/
# 	https://www.openpolicyagent.org/docs/latest/policy-reference/

# ==============================================================================
# Define dependencies

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
# Manage Project

up:
	encore run -v --browser never

upgrade:
	encore version update

resetdb:
	encore db reset app
	encore db reset test-app

reset-encore:
	cd "/Users/bill/Library/Application Support/encore"; \
	rm encore.db; \
	rm encore.db-shm; \
	rm encore.db-wal; \
	rm onboarding.json;

secrets:
	cat zarf/keys/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem | encore secret set --type local KeyPEM
	cat zarf/keys/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem | encore secret set --type dev KeyPEM
	echo "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1" | encore secret set --type local KeyID
	echo "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1" | encore secret set --type dev KeyID

metrics:
	expvarmon -ports="localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

statsviz:
	open -a "Google Chrome" http://127.0.0.1:4000/debug/statsviz

# ==============================================================================
# Shut Down

FIND_DB = $(shell docker ps | grep encoredotdev | cut -c 1-12)
SET_DB = $(eval DB_ID=$(FIND_DB))

FIND_DAEMON = $(shell ps | grep 'encore daemon' | grep -v 'grep' | cut -c 1-5)
SET_DAEMON = $(eval DAEMON_ID=$(FIND_DAEMON))

FIND_APP = $(shell ps | grep 'encore_app_out' | grep -v 'grep' | cut -c 1-5)
SET_APP = $(eval APP_ID=$(FIND_APP))

down-db:
	$(SET_DB)
	if [ -z "$(DB_ID)" ]; then \
		echo "db not running"; \
    else \
		docker stop $(DB_ID); \
		docker rm $(DB_ID) -v; \
    fi

down-daemon:
	$(SET_DAEMON)
	if [ -z "$(DAEMON_ID)" ]; then \
		echo "daemon not running"; \
    else \
		kill -SIGTERM $(DAEMON_ID); \
    fi

down-app:
	$(SET_APP)
	if [ -z "$(APP_ID)" ]; then \
		echo "app not running"; \
    else \
		kill -SIGTERM $(APP_ID); \
    fi

down: down-app down-daemon down-db

# ==============================================================================
# Running tests within the local computer

test-r:
	CGO_ENABLED=1 encore test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 encore test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test: test-only vuln-check lint

test-down: test-only vuln-check lint down

test-race: test-r vuln-check lint

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy

tidy:
	go mod tidy

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Administration

pgcli:
	pgcli $(shell encore db conn-uri app)

# ==============================================================================
# Hitting endpoints

token:
	curl -il \
	--user "admin@example.com:gophers" http://localhost:4000/v1/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

token-cloud:
	curl -il \
	--user "admin@example.com:gophers" http://app.encore.dev/sales-7a6i/v1/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

# export TOKEN="COPY TOKEN STRING FROM LAST CALL"

users:
	curl -il \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:4000/v1/users?page=1&rows=2"

load:
	hey -m GET -c 100 -n 1000 \
	-H "Authorization: Bearer ${TOKEN}" "http://localhost:4000/v1/users?page=1&rows=2"
