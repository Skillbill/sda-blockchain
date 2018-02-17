BINDIR?=bin
GO?=go

bin:=$(abspath $(BINDIR))

.PHONY:	all clean

all:	$(bin)/toychain $(bin)/hashcash
	@touch $(bin)

$(bin):
	mkdir -p $@

$(bin)/toychain: $(bin)
	cd src/cmd/toychain && $(GO) build -o $@

$(bin)/hashcash: $(bin)
	cd src/cmd && $(GO) build -o $@ $(notdir $@).go

clean:
	rm -f $(bin)/toychain $(bin)/hashcash
	test ! -e $(bin) || rmdir $(bin)
