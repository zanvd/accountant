# Dev build
FROM golang:1.17 AS dev-build

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux go build -gcflags="all=-N -l" -o ./dist/accountant .
RUN GOOS=linux go build -o ./dist/accountant-cmd ./cmd

# Dev container
FROM golang:1.17 AS dev

RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

COPY --from=dev-build /tmp/app/dist/accountant ./accountant
COPY --from=dev-build /tmp/app/dist/accountant-cmd ./accountant-cmd
COPY ./assets ./assets
COPY ./templates ./templates

CMD [ "dlv", "exec", "/app/accountant", "--headless=true", "--listen=:2345", "--api-version=2", "--accept-multiclient", "--log", "--continue" ]

# Prod build
FROM golang:1.17 AS prod-build

WORKDIR /tmp/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux go build -o ./dist/accountant .
RUN GOOS=linux go build -o ./dist/accountant-cmd ./cmd

# Prod container
FROM golang:1.17 AS prod

WORKDIR /app

COPY --from=prod-build /tmp/app/dist/accountant ./accountant
COPY --from=prod-build /tmp/app/dist/accountant-cmd ./accountant-cmd
COPY ./assets ./assets
COPY ./templates ./templates

CMD ["./accountant"]
