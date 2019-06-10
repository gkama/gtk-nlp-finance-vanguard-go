FROM golang:latest

WORKDIR /go/src/nlp
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

# Run the executable
CMD ["nlp"]