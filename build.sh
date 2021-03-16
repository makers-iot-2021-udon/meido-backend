#!/bin/bash
docker build ./flask --tag notchman/meido-flask:latest --build-arg flask_url=$1
docker push notchman/meido-flask:latest
