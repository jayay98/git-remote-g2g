.PHONY: all clean install build fmt

objects := git-g2g git-remote-g2g

all: fmt build

clean:
	$(foreach obj, $(objects), go clean -C cmd/$(obj) -i;)

install:
	$(foreach obj, $(objects), go install -C cmd/$(obj);)

build:
	$(foreach obj, $(objects), go build -C cmd/$(obj) || exit 1;)

fmt:
	gofmt -w -l .