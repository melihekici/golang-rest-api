FROM golang:latest
RUN mkdir /app
COPY . /usr/src/app
WORKDIR /usr/src/app
RUN go build -o main /usr/src/app
CMD ["/usr/src/app/main"]
