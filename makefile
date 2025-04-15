libcryptage_version=v0.0.2
ageversion = $(shell cat go.mod | grep io/age | awk '{print $$3}')
xcryptoversion = $(shell cat go.mod | grep x/crypto | awk '{print $$2}')
goldflags += -X main.version=${libcryptage_version}
goldflags += -X main.ageversion=${ageversion}
goldflags += -X main.xcryptoversion=${xcryptoversion}
goflags=-ldflags '${goldflags}'
library: libcryptage.a libcryptage.so
	@echo "run 'make all' to build the examples too"
help:
	@echo run:
	@echo "	make clean all"
	@echo "	sudo make install"
example/example:
	${MAKE} -C example
all: libcryptage.a libcryptage.so example/example
	@echo "run 'sudo make install' to install the library system-wide"
libcryptage.a: *.go
	go build -v ${goflags} -o $@ -buildmode c-archive -tags 'netgo,osusergo' .
	@sha256sum $@
	@file $@
libcryptage.so: *.go
	go build -v ${goflags} -o $@ -buildmode c-shared -tags 'netgo,osusergo' .
	@sha256sum $@
	@file $@
clean:
	rm -f *.a *.so
	${MAKE} -C example clean
install: # no dep prevent sudo install from building go files as root
	@test -f libcryptage.a || echo "first, run make"
	@test -f libcryptage.a
	@test -f libcryptage.so
	mv -v libcryptage.a libcryptage.so ${DESTDIR}/usr/local/lib/
	@test -f libcryptage.h
	cp -v libcryptage.h ${DESTDIR}/usr/local/include/
	@echo "run 'make -C example postinstall' to try building the example using system includes dir"

