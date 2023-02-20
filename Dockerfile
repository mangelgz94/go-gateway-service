FROM alpine:latest as alpine
RUN apk add --no-cache tzdata ca-certificates

FROM golang:1.19 as gobuild

ENV CGO_ENABLED=0

WORKDIR /opt
COPY . .
RUN go build -o build/users -ldflags "-X main.BuildVersion=1" ./cmd/users
RUN go build -o build/find_number_position -ldflags "-X main.BuildVersion=1" ./cmd/find_number_position
RUN go build -o build/gateway -ldflags "-X main.BuildVersion=1" ./cmd/gateway

RUN groupadd -g 3000 appuser && useradd -r -u 1000 -g appuser appuser
RUN chown -R appuser:appuser build/users
RUN chown -R appuser:appuser build/find_number_position
RUN chown -R appuser:appuser build/gateway

FROM scratch as release

COPY --from=alpine /etc/ssl/certs /etc/ssl/certs

COPY --from=gobuild /etc/passwd /etc/passwd
COPY --from=gobuild /opt/build/users opt/users
COPY  users_json/ /opt/users_json
COPY --from=gobuild /opt/build/find_number_position opt/find_number_position
COPY --from=gobuild /opt/build/gateway opt/gateway

USER appuser
WORKDIR /opt
ENTRYPOINT ["/opt/gateway"]
EXPOSE 8090
EXPOSE 50051
EXPOSE 50052

