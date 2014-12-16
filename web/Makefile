PID      = /tmp/logvoyage-frontend.pid
GO_FILES = $(wildcard *.go)
 
serve:
	@make restart 
	@fswatch -o . | xargs -n1 -I{}  make restart || make kill
 
kill:
	@kill `cat $(PID)` || true
 
restart:
	echo "Restarting..."
	@make kill
	@go run $(GO_FILES) & echo $$! > $(PID)
 
.PHONY: serve restart kill 