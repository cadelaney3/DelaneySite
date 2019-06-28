FROM golang:1.11

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

# RUN go mod download

COPY . .

# RUN go get -d -v ./...

# RUN go install -v ./...

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 8080

CMD ["delaneySite"]
