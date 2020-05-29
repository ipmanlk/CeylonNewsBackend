#!/bin/bash
docker build --tag ceylon_news .
docker stop cn
docker rm cn
# docker run -d -p 3000:3000 --name cn ceylon_news 

docker run -d \
  -it \
  -p 3000:3000 \
  --name cn \
  --mount type=bind,source="$(pwd)"/auth,target=/usr/src/app/src/api/auth \
  ceylon_news
