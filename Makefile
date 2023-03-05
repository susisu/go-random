.PHONY: test
test:
	go test -v -race ./...

.PHONY: update-snapshots
update-snapshots:
	UPDATE_SNAPS=true $(MAKE) test
