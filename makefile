library: libcryptage.a libcryptage.so
libcryptage.a: *.go
	go build -v -o $@ -buildmode c-archive -tags 'netgo,osusergo' *.go
	@sha256sum $@
	@file $@
libcryptage.so: *.go
	go build -v -o $@ -buildmode c-shared -tags 'netgo,osusergo' *.go
	@sha256sum $@
	@file $@
clean:
	rm -f *.a *.so
