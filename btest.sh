go build  ./src/gtest/gtest.go
rm -rf producttest
mkdir producttest
mv gtest ./producttest/testServer
producttest/testServer
