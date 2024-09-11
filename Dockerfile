FROM golang as builder

RUN apt update && apt install -y exiftool

WORKDIR /usr/src

COPY go.mod go.sum ./
RUN go mod download && go mod verify && go mod tidy

COPY . .

RUN go build -v -o app ./...

FROM golang

WORKDIR /usr/local/bin

ENV GIN_MODE=release

COPY --from=builder /usr/src/ /usr/local/bin/

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/app"]