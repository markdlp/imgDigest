FROM golang

WORKDIR /usr/src/app

RUN apt update && apt install exiftool -y

ENV GIN_MODE=release

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

RUN go install golang.org/x/tools/gopls@latest \
&& go install github.com/cweill/gotests/gotests@latest \
&& go install github.com/fatih/gomodifytags@latest \
&& go install github.com/josharian/impl@latest \
&& go install github.com/haya14busa/goplay/cmd/goplay@latest \
&& go install github.com/go-delve/delve/cmd/dlv@latest \
&& go install honnef.co/go/tools/cmd/staticcheck@latest

EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/app" ]
