# Build container
FROM golang:1.15 AS build

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux go build -o ./dist/accountant .

# Running container
FROM golang:1.15

COPY --from=build /tmp/app/dist/accountant /app/accountant
COPY ./assets /app/assets
COPY ./templates /app/templates

WORKDIR /app

EXPOSE 8080

CMD ["./accountant"]
