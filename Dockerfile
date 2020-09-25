# pull go alpine and install git
FROM golang:alpine
RUN apk update && apk add --no-cache git

# build apps
WORKDIR /app
ADD . /app
RUN go mod download
RUN go mod verify
RUN go build -o todo .

#CMD ["/app/todo"]
