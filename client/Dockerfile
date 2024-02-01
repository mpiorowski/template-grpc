# Development
FROM node:20-slim

WORKDIR /client

COPY package.json /client/package.json
COPY package-lock.json /client/package-lock.json

ARG ENV
RUN echo "ENV=${ENV}" >> .env

RUN npm install

COPY . .

CMD npm run dev
