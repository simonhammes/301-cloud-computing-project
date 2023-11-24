FROM golang:1.21-alpine AS build

RUN apk add --update git hugo

WORKDIR /app

COPY slides .

RUN hugo

FROM nginx:1.25

WORKDIR /usr/share/nginx/html

COPY --from=build /app/public /usr/share/nginx/html
