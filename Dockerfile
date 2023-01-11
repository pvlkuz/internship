FROM golang:alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

#RUN migrate -path ./ -database 'postgres://postgres:password@localhost:5432/postgres?sslmode=disable' up
RUN go build -o ./bin/main

EXPOSE 8080

CMD [ "./bin/main" ]