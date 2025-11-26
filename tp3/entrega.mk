netstats:
	go mod init tp3 2>/dev/null || true
	go build -o netstats

clean:
	rm -f netstats
