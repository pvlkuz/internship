FROM golang

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o ./bin/main

EXPOSE 8080

#CMD [ "./bin/main" ]
ENTRYPOINT [ "./bin/main", "crud" ]