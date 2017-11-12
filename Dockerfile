FROM golang:1.8.5-jessie

COPY src /src

CMD ["go", "run", "/src/wiki.go"]
