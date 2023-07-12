PREFIX = /usr/local
COMPLETIONS_DIR_BASH = $(PREFIX)/share/bash-completion/completions
COMPLETIONS_DIR_ZSH = $(PREFIX)/share/zsh/site-functions
COMPLETIONS_DIR_FISH = $(PREFIX)/share/fish/vendor_completions.d

all: przypominajka completions

przypominajka:
	 go build -ldflags "-X main.version=$$(git describe --always --dirty)" .

completions: przypominajka.bash przypominajka.zsh przypominajka.fish

przypominajka.bash: przypominajka
	./przypominajka completion bash > przypominajka.bash

przypominajka.zsh: przypominajka
	./przypominajka completion zsh > przypominajka.zsh

przypominajka.fish: przypominajka
	./przypominajka completion fish > przypominajka.fish

clean: 
	rm -f przypominajka przypominajka.bash przypominajka.zsh przypominajka.fish

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

uninstall:
	rm -f \
		$(PREFIX)/bin/przypominajka \
		$(COMPLETIONS_DIR_BASH)/przypominajka \
		$(COMPLETIONS_DIR_ZSH)/_przypominajka \
		$(COMPLETIONS_DIR_FISH)/przypominajka.fish

.PHONY: all przypominajka completions clean install uninstall
