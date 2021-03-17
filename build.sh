#!/bin/bash

echo $1
docker build ./flask --tag notchman/meido-flask:latest

docker push notchman/meido-flask:latest
