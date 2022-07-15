FROM golang:1.17.7-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /docker-server

##
## Deploy
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /docker-server /docker-server
COPY --chown=nonroot:nonroot lookup_table.yml ./
EXPOSE 8087

USER nonroot:nonroot

ENTRYPOINT ["/docker-server"]