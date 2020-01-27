FROM golang:latest
RUN mkdir /app
WORKDIR /app
COPY ./go.mod .
RUN go mod download
COPY ./ /app
RUN GOOS=linux GOARCH=amd64 go build
CMD /app/worker
