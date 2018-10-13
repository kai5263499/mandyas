# build stage
FROM golang:alpine AS build

ENV GOPATH=/app

WORKDIR /app

RUN apk update && \
    apk add git && \
    go get -u github.com/golang/dep/cmd/dep && \
    go get -u github.com/golang/protobuf/protoc-gen-go && \
    cp /app/bin/dep /usr/local/bin && \
    cp /app/bin/protoc-gen-go /usr/local/bin

COPY src/mandyas-service /app/src/mandyas-service
RUN cd /app/src/mandyas-service && \
    make

FROM itzg/minecraft-server

COPY --from=build /app/src/mandyas-service/mandyas /mandyas

ENV EULA=TRUE

ENTRYPOINT /mandyas /start
