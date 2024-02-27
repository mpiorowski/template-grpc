# Development
FROM node:20-slim AS dev

WORKDIR /client

COPY package.json /client/package.json
COPY package-lock.json /client/package-lock.json

ARG TARGET
RUN echo "TARGET=${TARGET}" >> .env
ARG JWT_SECRET
RUN echo "JWT_SECRET=${JWT_SECRET}" >> .env
ARG PUBLIC_AUTH_URL 
RUN echo "PUBLIC_AUTH_URL=${PUBLIC_AUTH_URL}" >> .env
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
RUN npm install
COPY . .

ARG TARGET
RUN echo "TARGET=${TARGET}" >> .env
ARG JWT_SECRET
RUN echo "JWT_SECRET=${JWT_SECRET}" >> .env
ARG PUBLIC_AUTH_URL 
RUN echo "PUBLIC_AUTH_URL=${PUBLIC_AUTH_URL}" >> .env
ARG DIRECTUS_URL
RUN echo "DIRECTUS_URL=${DIRECTUS_URL}" >> .env
ARG USERS_URI
RUN echo "USERS_URI=${USERS_URI}" >> .env

RUN npm run build

# Production
FROM node:20-slim AS prod
WORKDIR /client

COPY package.json /client/package.json
COPY package-lock.json /client/package-lock.json
RUN npm install --omit=dev

COPY --from=build /client/build /client/build
COPY --from=build /client/src/lib/proto/ /client/src/lib/proto/

CMD PORT=3000 node build

