FROM alpine:latest

FROM golang:1.18-alpine

WORKDIR /home 

COPY labCode/ .

RUN apk update && apk add curl

RUN go mod download

RUN go build -o output *.go

# CMD ["./output"]

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#

# $ docker build . -t kadlab
