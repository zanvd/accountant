# Build container
FROM golang:1.15 AS build

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux go build -o ./dist/accountant .
RUN GOOS=linux go build -o ./dist/accountant-cmd ./cmd

# Running container
FROM golang:1.15

COPY --from=build /tmp/app/dist/accountant /app/accountant
COPY --from=build /tmp/app/dist/accountant-cmd /app/accountant-cmd
COPY ./assets /app/assets
COPY ./templates /app/templates

WORKDIR /app

EXPOSE 8080

CMD ["./accountant"]
