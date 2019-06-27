go build -a./src/LB/LBmain.go 
mv LBmain ./product/LBServer
go build -a ./src/GATE/GTmain.go
mv GTmain ./product/GTserver
