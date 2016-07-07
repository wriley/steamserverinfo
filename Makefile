steamserverinfo: steamserverinfo.go
	@go build

install: steamserverinfo
	@install -m 0755 steamserverinfo /usr/local/bin/

clean:
	@/bin/rm -f steamserverinfo
