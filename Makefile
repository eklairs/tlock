# Installs all the git hooks
install-git-hooks:
	cp hooks/* .git/hooks/

# Generates vendor
generate-vendor:
	go run genvendor.go

# Builds git version of the tlock
# Basically sets the version to the latest commit
build-git:
	@git_version=$$(git rev-list --count HEAD)"+"$$(git rev-parse --short HEAD); \
		go build -ldflags "-X github.com/eklairs/tlock/tlock-internal/constants.VERSION=v$$git_version" -o tlock-v$$git_version tlock/main.go

