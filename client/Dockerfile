# Development
FROM node:20-slim AS dev

WORKDIR /client

COPY package.json /client/package.json
COPY package-lock.json /client/package-lock.json

ARG TARGET
RUN echo "TARGET=${TARGET}" >> .env
ARG JWT_SECRET
RUN echo "JWT_SECRET=${JWT_SECRET}" >> .env
ARG PUBLIC_AUTH_HTTP 
RUN echo "PUBLIC_AUTH_HTTP=${PUBLIC_AUTH_HTTP}" >> .env
ARG DIRECTUS_URL
RUN echo "DIRECTUS_URL=${DIRECTUS_URL}" >> .env
ARG USERS_URI
RUN echo "USERS_URI=${USERS_URI}" >> .env

RUN npm install

COPY . .

# Build
FROM node:20-slim AS build
WORKDIR /client

COPY package.json /client/package.json
COPY package-lock.json /client/package-lock.json

ARG TARGET
RUN echo "TARGET=${TARGET}" >> .env
ARG JWT_SECRET
RUN echo "JWT_SECRET=${JWT_SECRET}" >> .env
ARG PUBLIC_AUTH_HTTP 
RUN echo "PUBLIC_AUTH_HTTP=${PUBLIC_AUTH_HTTP}" >> .env
ARG DIRECTUS_URL
RUN echo "DIRECTUS_URL=${DIRECTUS_URL}" >> .env
ARG USERS_URI
RUN echo "USERS_URI=${USERS_URI}" >> .env

RUN npm install
COPY . .

RUN npm run build

# Production
FROM node:20-slim AS prod
WORKDIR /client

COPY package.json /client/package.json
COPY package-lock.json /client/package-lock.json
RUN npm install --production

COPY --from=build /client/build /client/build

CMD PORT=3000 node build

