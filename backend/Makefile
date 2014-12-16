PID      = /tmp/logvoyage-backend.pid
GO_FILES = $(wildcard *.go)
 
serve:
	@make restart 
	@fswatch -o . | xargs -n1 -I{}  make restart || make kill
 
kill:
	@kill `cat $(PID)` || true
 
restart:
	@make kill
	@go run $(GO_FILES) & echo $$! > $(PID)
 
.PHONY: serve restart kill 