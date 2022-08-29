# get a base image
FROM golang:1.17-alpine

# set the working directory at the container
WORKDIR /go/src/app

# copy the files from host to the container working directory
COPY ./app ./

# downlod all the dependecies listed in the go.mod
RUN go get -d -v

# build the project into a binary
RUN go build -v

# run the binary after build
CMD ["./golang-mongodb-api"]
