npm run build
~/go/bin/go-bindata build/...
sed -i 's/package main/package web/g' bindata.go
mv bindata.go ../go/web/
rm -rf ./build
