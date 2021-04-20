#!/bin/bash
echo Building the Todo CLI...
sudo mv ./todo /usr/local 
cd /usr/local/todo
go build todo.go
mv ./todo /usr/local/bin

echo Building has finished


