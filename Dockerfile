# Dockerfile
FROM golang:1.4-onbuild
EXPOSE 8080
CMD make && ./bin/crawler
