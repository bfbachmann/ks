.PHONY: install
install:
	go build && go install

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor: tidy
	go mod vendor

.PHONY: release
release:
	gox -output="bin/ks_{{.OS}}_{{.Arch}}"
