# Development
FROM node:20-slim AS dev

WORKDIR /client

COPY package.json /client/package.json
COPY package-lock.json /client/package-lock.json
COPY .npmrc /client/.npmrc

ARG TARGET
RUN echo "TARGET=${TARGET}" >> .env
ARG JWT_SECRET
RUN echo "JWT_SECRET=${JWT_SECRET}" >> .env
ARG PUBLIC_AUTH_URL 
RUN echo "PUBLIC_AUTH_URL=${PUBLIC_AUTH_URL}" >> .env
ARG AUTH_URI
RUN echo "AUTH_URI=${AUTH_URI}" >> .env
ARG PROFILE_URI
RUN echo "PROFILE_URI=${PROFILE_URI}" >> .env

RUN npm install

COPY . /client

# Build
FROM node:20-slim AS build
WORKDIR /client

COPY package.json /client/package.json
COPY package-lock.json /client/package-lock.json
COPY .npmrc /client/.npmrc

RUN npm install

RUN --mount=type=secret,id=JWT_SECRET echo "JWT_SECRET=$(cat /run/secrets/JWT_SECRET)" >> .env

ARG TARGET
RUN echo "TARGET=${TARGET}" >> .env
# ARG JWT_SECRET
# RUN echo "JWT_SECRET=${JWT_SECRET}" >> .env
ARG PUBLIC_AUTH_URL 
RUN echo "PUBLIC_AUTH_URL=${PUBLIC_AUTH_URL}" >> .env
ARG AUTH_URI
RUN echo "AUTH_URI=${AUTH_URI}" >> .env
ARG PROFILE_URI
RUN echo "PROFILE_URI=${PROFILE_URI}" >> .env

COPY . /client
RUN npm run build

# Production
FROM node:20-slim AS prod
WORKDIR /client

COPY package.json /client/package.json
COPY package-lock.json /client/package-lock.json
COPY .npmrc /client/.npmrc
RUN npm install --omit=dev

COPY --from=build /client/build /client/build
COPY --from=build /client/src/lib/proto/ /client/src/lib/proto/

CMD PORT=3000 node build

