# Build react app
FROM node:14.15.4-alpine as client_builder

RUN apk add --no-cache --update git

WORKDIR /opt/poker-app-client
COPY ./client/poker-app/package.json .
COPY ./client/poker-app/yarn.lock .
RUN yarn install

COPY ./client/poker-app .
RUN yarn build

# Build game server
FROM gcr.io/gcp-runtimes/go1-builder:1.15 as server_builder

COPY . /go/src/poker-app

WORKDIR /go/src/poker-app/cmd/poker-app
RUN /usr/local/go/bin/go build

# Application image
FROM gcr.io/distroless/base:latest

COPY --from=server_builder /go/src/poker-app/cmd/poker-app/poker-app /usr/local/bin/poker-app
COPY --from=client_builder /opt/poker-app-client/build /usr/local/lib/poker-app/client

CMD ["/usr/local/bin/poker-app"]
