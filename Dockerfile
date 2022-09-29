FROM golang:1.16-alpine

WORKDIR app/

COPY . .

RUN go mod download && go build cmd/build.go

ARG PORT=8080

ENV PORT ${PORT}

EXPOSE ${PORT}

CMD ["sh", "-c", "./build -p ${PORT}"]