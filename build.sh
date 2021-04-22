#/usr/bin/sh

#Web Build
cd web
./build.sh

#Go Build
cd ../go
export GOX_windows_amd64_LDFLAGS="-H=windowsgui"
~/go/bin/gox -osarch="linux/amd64 linux/arm windows/amd64" -output="../sdist/bin/ash_{{.OS}}_{{.Arch}}"

cd ../sdist
mkdir ../dist

#Debian Package
chmod -R 0755 debian
chmod 0644 debian/almost-scrum-0.5/usr/share/doc/almost-scrum/*
strip -s -o debian/almost-scrum-0.5/usr/bin/ash bin/ash_linux 
fakeroot dpkg-deb --build debian/almost-scrum-0.5
mv debian/almost-scrum-0.5.deb ../dist

#Windows Package
cd windows
makensis ash-setup.nsi
cd ..
