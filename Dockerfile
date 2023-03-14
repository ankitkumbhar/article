# base image
FROM golang:1.19-alpine as builder

# create app dir
RUN mkdir /app

# add all files to app dir
ADD . /app

# set workdir to app
WORKDIR /app/cmd

# build the application
RUN go build -o main .

# final stage
FROM alpine:latest

# create app dir
RUN mkdir /app

# set workdir to app
WORKDIR /app

# copy prebuilt binary from previous stage to new stage
COPY --from=builder /app/cmd/main /app

# expose port to outside world
EXPOSE 8080

# run the executable
ENTRYPOINT ./main
