#!/usr/bin/env bash

go test -c -o repository.test github.com/dacero/labyrinth-of-babel/repository && ./repository.test
go test -c -o handlers.test github.com/dacero/labyrinth-of-babel/handlers && ./handlers.test

#go test github.com/dacero/labyrinth-of-babel/repository
#go test github.com/dacero/labyrinth-of-babel/handlers