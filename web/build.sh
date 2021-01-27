npm run build
~/go/bin/go-bindata build/... assets/...
sed -i 's/package main/package assets/g' bindata.go
mv bindata.go ../go/assets/
rm -rf ./build
