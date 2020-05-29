#!/bin/bash
docker pull ipmanlk/ceylonnews
docker stop ceylon_news
docker rm ceylon_news
docker run -d \
  -it \
  -p 3000:3000 \
  --name ceylon_news \
  --restart always \
  --mount type=bind,source="$(pwd)"/auth,target=/usr/src/app/src/api/auth \
  ipmanlk/ceylonnews