FROM golang:1.19

WORKDIR /opt/validator
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 80

CMD ["sh", "-c", "validator"]