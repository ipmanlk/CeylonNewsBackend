#!/bin/bash
docker build --tag ceylon_news .
docker stop ceylon_news
docker rm ceylon_news
# docker run -d -p 3000:3000 --name cn ceylon_news 

docker run -d \
  -it \
  -p 3000:3000 \
  --name ceylon_news \
  --mount type=bind,source="$(pwd)"/.env,target=/usr/src/app/.env \
  --mount type=bind,source="$(pwd)"/auth,target=/usr/src/app/src/api/auth \
  ceylon_news
