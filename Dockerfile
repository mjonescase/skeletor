# create image from the official Go image
FROM golang

# Use Go Modules instead of GOPATH and a 3rd party dep. manager
ENV GO111MODULE=on
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# install fresh
# RUN go get github.com/pilu/fresh

EXPOSE 5000

# serve the app
ENTRYPOINT ["/app/skeletor"]