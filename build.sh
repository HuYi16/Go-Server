go build  ./src/LB/LBmain.go 
mv LBmain ./product/LBServer
go build  ./src/GATE/GTmain.go
mv GTmain ./product/GTserver
