#! /bin/bash

for i in $(find . -iname "*_test.go" -exec dirname {} \; | uniq)
do
    go test -cover $i;
    if [ $? != 0 ]
    then
      exit 1
    fi
done