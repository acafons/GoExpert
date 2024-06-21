FROM golang:latest

WORKDIR /usr/src/app

RUN apt update -y && apt install -y sqlite3

EXPOSE 8080
