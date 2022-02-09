# Build container
FROM golang:1.17 AS build

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux go build -gcflags="all=-N -l" -o ./dist/accountant .
RUN GOOS=linux go build -o ./dist/accountant-cmd ./cmd

# Running container
FROM golang:1.17

# TODO Install this for dev only.
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

COPY --from=build /tmp/app/dist/accountant ./accountant
COPY --from=build /tmp/app/dist/accountant-cmd ./accountant-cmd
COPY ./assets ./assets
COPY ./templates ./templates

#CMD ["./accountant"]
CMD [ "dlv", "exec", "/app/accountant", "--headless=true", "--listen=:2345", "--api-version=2", "--accept-multiclient", "--log", "--continue" ]
