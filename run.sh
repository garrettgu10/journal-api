#!/bin/bash

export $(xargs <.env)
mkdir -p tmp
go run *.go