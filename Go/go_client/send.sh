#!/bin/bash

nc localhost 5828 < com-amazon.txt > communities.txt
go run main.go
cat communitiesTitle.txt
