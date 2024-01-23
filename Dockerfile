FROM node:16-buster

WORKDIR /usr/src/app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run ready

RUN npm run build

EXPOSE 5000

CMD ["node", "dist/app.js"]
