go build  ./src/LB/LBmain.go
rm -rf product
mkdir product
mv LBmain ./product/LBServer
product/LBServer
