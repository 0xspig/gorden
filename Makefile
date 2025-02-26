PREFIX = /usr/local

build: ./app/* ./garden/* ./server/*
	go build -o ./build/gorden ./app/gorden.go
install: build
	#TODO copy def files into appropriate system directory
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f ./build/gorden ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/gorden
	mkdir -p ${DESTDIR}${PREFIX}/share/gorden
	cp -rf ./share/* ${DESTDIR}${PREFIX}/share/gorden
uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/gorden
	rm -f ${DESTDIR}${PREFIX}/share/gorden
