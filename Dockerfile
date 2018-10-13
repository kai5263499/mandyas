# build stage
FROM golang:alpine

WORKDIR /app
ADD . /app
RUN cd /app && go build -o goapp

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/goapp /app/
ENTRYPOINT ./goapp
