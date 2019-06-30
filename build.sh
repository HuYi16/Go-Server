go build  ./src/LB/LBmain.go
rm -rf product
mkdir product
mv LBmain ./product/LBServer
product/LBServer
go build  ./src/GATE/GTmain.go
mv GTmain ./product/GTserver
