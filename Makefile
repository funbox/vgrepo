########################################################################################

.PHONY = fmt all clean deps

########################################################################################

all: vgrepo

vgrepo:
	go build vgrepo.go

deps:
	git config --global http.https://gopkg.in.followRedirects true
	git config --global http.https://pkg.re.followRedirects true
	go get -d -v pkg.re/essentialkaos/ek.v8

fmt:
	find . -name "*.go" -exec gofmt -s -w {} \;

clean:
	rm -f vgrepo

########################################################################################
