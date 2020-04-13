#!/bin/bash
docker build --tag ceylon_news .
docker stop cn
docker rm cn
docker run -d -p 3000:3000 --name cn ceylon_news 