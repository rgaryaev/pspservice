FROM golang:latest

WORKDIR /pspservice/src

COPY . .

RUN go build -o pspservice .


WORKDIR /pspservice/app
RUN mkdir ./.data

RUN cp /pspservice/src/config.json .
RUN mv /pspservice/src/pspservice .
 
EXPOSE 8080

CMD ["./pspservice"]
