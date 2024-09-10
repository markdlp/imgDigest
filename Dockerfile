FROM golang

WORKDIR /usr/src/app

RUN apt update && apt install -y exiftool

ENV GIN_MODE=release

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

EXPOSE 8080
# ENTRYPOINT [ "app" ]
CMD [ "app" ]