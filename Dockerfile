FROM golang:latest

RUN mkdir /appfinal
ADD . /appfinal
WORKDIR /appfinal
RUN go build -o main .

CMD [ "/appfinal/main" ]

