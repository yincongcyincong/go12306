CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go12306-linux
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o go12306-windows.exe
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o go12306-mac

tar -zcvf go12306-linux.tar.gz ./go12306-linux ./conf/conf.ini
tar -zcvf go12306-windows.tar.gz ./go12306-windows.exe ./conf/conf.ini
tar -zcvf go12306-mac.tar.gz ./go12306-mac ./conf/conf.ini
rm -rf go12306-linux go12306-windows.exe go12306-mac