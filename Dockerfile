FROM ubuntu:latest
LABEL authors="Yuxing He"

WORKDIR $GOPATH/src/nova
COPY . .

RUN go get -v -t -d ./...
RUN go install -v ./...

EXPOSE 8080
ENTRYPOINT ["$GOPATH/bin/nova"]