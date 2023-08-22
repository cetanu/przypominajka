SHELL = /bin/sh

PREFIX = /usr/local
COMPLETIONS_DIR_BASH = $(PREFIX)/share/bash-completion/completions
COMPLETIONS_DIR_ZSH = $(PREFIX)/share/zsh/site-functions
COMPLETIONS_DIR_FISH = $(PREFIX)/share/fish/vendor_completions.d


.PHONY: all
all: przypominajka completions

.PHONY: przypominajka
przypominajka:
	 go build -ldflags "-X main.version=$$(git describe --always --dirty)" .

.PHONY: completions
completions: przypominajka.bash przypominajka.zsh przypominajka.fish

.PHONY: przypominajka.bash
przypominajka.bash: przypominajka
	./przypominajka completion bash > przypominajka.bash

.PHONY: przypominajka.zsh
przypominajka.zsh: przypominajka
	./przypominajka completion zsh > przypominajka.zsh

.PHONY: przypominajka.fish
przypominajka.fish: przypominajka
	./przypominajka completion fish > przypominajka.fish

.PHONY: install
install:
	install -d \
		$(PREFIX)/bin \
		$(COMPLETIONS_DIR_BASH) \
		$(COMPLETIONS_DIR_ZSH) \
		$(COMPLETIONS_DIR_FISH)

	install -pm 0755 przypominajka $(PREFIX)/bin/przypominajka
	install -pm 0644 przypominajka.bash $(COMPLETIONS_DIR_BASH)/przypominajka
	install -pm 0644 przypominajka.zsh $(COMPLETIONS_DIR_ZSH)/_przypominajka
	install -pm 0644 przypominajka.fish $(COMPLETIONS_DIR_FISH)/przypominajka.fish

.PHONY: uninstall
uninstall:
	rm -f \
		$(PREFIX)/bin/przypominajka \
		$(COMPLETIONS_DIR_BASH)/przypominajka \
		$(COMPLETIONS_DIR_ZSH)/_przypominajka \
		$(COMPLETIONS_DIR_FISH)/przypominajka.fish

.PHONY: clean
clean: 
	rm -f przypominajka przypominajka.bash przypominajka.zsh przypominajka.fish
