# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.26-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -o /out/openNAMU .

FROM docker.io/library/debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates libgcc-s1 \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /data && chown 10001:10001 /data
COPY --from=build --chown=10001:10001 /out/openNAMU /usr/local/bin/openNAMU

ENV HOME=/data
WORKDIR /data
VOLUME ["/data"]
EXPOSE 3000
USER 10001:10001
ENTRYPOINT ["/usr/local/bin/openNAMU"]
