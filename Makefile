# Installs all the git hooks
install-git-hooks:
	cp hooks/* .git/hooks/

# Generates vendor
generate-vendor:
	go run genvendor.go

