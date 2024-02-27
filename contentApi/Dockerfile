FROM golang:1.22

WORKDIR /app
COPY . /app

COPY go.mod go.sum ./
RUN go mod download
RUN go build -o /main
CMD /main