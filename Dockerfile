FROM golang:latest as builder

LABEL maintaner="Adam Czerepinski <aczerepinski@gmail.com>"

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app .

EXPOSE 8080

CMD ["./main"] 
