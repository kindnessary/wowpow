FROM golang:alpine as build

WORKDIR /wowpow

COPY cmd/wowserver cmd/wowserver
COPY services/wowserver services/wowserver
COPY services/pkg services/pkg
COPY go.mod go.sum ./

RUN go mod download
RUN cd cmd/wowserver && go build -o /bin/server

FROM alpine:latest

COPY --from=build /bin/server /server

CMD /server
