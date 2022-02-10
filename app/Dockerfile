# Build
FROM golang:1.17 AS build

ARG GOARG=""

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN /bin/bash -c "GOOS=linux go build ${GOARG} -o ./dist/accountant ."
RUN GOOS=linux go build -o ./dist/accountant-cmd ./cmd

# Dev container
FROM golang:1.17 AS dev

RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

COPY --from=build /tmp/app/dist/accountant ./accountant
COPY --from=build /tmp/app/dist/accountant-cmd ./accountant-cmd
COPY ./assets ./assets
COPY ./templates ./templates

CMD [ "dlv", "exec", "/app/accountant", "--headless=true", "--listen=:2345", "--api-version=2", "--accept-multiclient", "--log", "--continue" ]

# Prod container
FROM golang:1.17 AS prod

WORKDIR /app

COPY --from=build /tmp/app/dist/accountant ./accountant
COPY --from=build /tmp/app/dist/accountant-cmd ./accountant-cmd
COPY ./assets ./assets
COPY ./templates ./templates

CMD ["./accountant"]
