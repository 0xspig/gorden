PREFIX = /usr/local

build:
	go build -o ./build/gorden ./app/gorden.go
install: build
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f ./build/gorden ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/gorden
uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/gorden
