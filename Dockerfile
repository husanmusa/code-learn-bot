FROM golang:1.17-alpine

RUN mkdir /application

COPY . /application

WORKDIR /application
RUN go mod tidy -go=1.16 && go mod tidy -go=1.17

RUN go run main.go

CMD ./main