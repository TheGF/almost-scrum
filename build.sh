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

#Debian Package amd64
chmod -R 0755 debian
chmod 0644 debian/almost-scrum_0.5_amd64/usr/share/doc/almost-scrum/*
strip -s -o debian/almost-scrum_0.5_amd64/usr/bin/ash bin/ash_linux_amd64 
fakeroot dpkg-deb --build debian/almost-scrum_0.5_amd64
mv debian/almost-scrum_0.5_amd64.deb ../dist

#Debian Package arm
chmod 0644 debian/almost-scrum_0.5_arm/usr/share/doc/almost-scrum/*
strip -s -o debian/almost-scrum_0.5_arm/usr/bin/ash bin/ash_linux_arm
fakeroot dpkg-deb --build debian/almost-scrum_0.5_arm
mv debian/almost-scrum_0.5_arm.deb ../dist


#Windows Package
cd windows
makensis ash-setup.nsi
cd ..
