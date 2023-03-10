FROM golang:alpine as build

WORKDIR /wowpow

COPY cmd/wowclient cmd/wowclient
COPY services/wowclient services/wowclient
COPY services/pkg services/pkg
COPY go.mod go.sum ./

RUN go mod download
RUN cd cmd/wowclient && go build -o /bin/server

FROM alpine:latest

COPY --from=build /bin/server /server

CMD /server
