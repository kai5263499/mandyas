FROM golang:alpine AS builder

ENV GOPATH=/app

RUN mkdir -p /app/src/github.com/kai5263499/mandyas && \
    apk update && \
    apk -U add git make curl protobuf && \
    go get -u github.com/golang/dep/cmd/dep && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    cp /app/bin/dep /usr/local/bin/dep && \
    cp /app/bin/protoc-gen-go /usr/local/bin/protoc-gen-go

WORKDIR /app/src/github.com/kai5263499/mandyas

COPY . /app/src/github.com/kai5263499/mandyas
RUN cd /app/src/github.com/kai5263499/mandyas && make

FROM itzg/minecraft-server

ENV EULA=TRUE

COPY --from=builder /app/src/github.com/kai5263499/mandyas/mandyas /mandyas

EXPOSE 25565 25575 9000

ENTRYPOINT ["/mandyas", "/start"]
