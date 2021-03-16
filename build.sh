#!/bin/bash

echo $1
docker build ./meido-api --tag notchman/meido-flask:latest --build-arg flask_url=$1

docker push notchman/meido-flask:latest
