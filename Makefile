.PHONY: all clean install build

objects := git-g2g git-remote-g2g

all: build

clean:
	$(foreach obj, $(objects), go clean -C cmd/$(obj) -i;)

install:
	$(foreach obj, $(objects), go install -C cmd/$(obj);)

build:
	$(foreach obj, $(objects), go build -C cmd/$(obj) || exit 1;)