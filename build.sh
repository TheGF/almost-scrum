#/usr/bin/sh

#Web Build
cd web
./build.sh

#Go Build
cd ../go
~/go/bin/gox -osarch="linux/amd64 windows/amd64 darwin/amd64" -output="../dist/bin/ash_{{.OS}}"

cd ../dist

#Debian Package
chmod -R 0755 debian
chmod 0644 debian/almost-scrum-0.5/usr/share/doc/almost-scrum/*
strip -s -o debian/almost-scrum-0.5/usr/bin/ash bin/ash_linux 
fakeroot dpkg-deb --build debian/almost-scrum-0.5
