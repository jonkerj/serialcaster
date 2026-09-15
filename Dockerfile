# ── build stage ───────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN GOPROXY=https://nexus.equinix.com/repository/golang-proxy/ \
    go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /serialcaster .

# ── runtime stage ─────────────────────────────────────────────────────────────
FROM scratch

# Override at build time: docker build --build-arg DIALOUT_GID=$(getent group dialout | cut -d: -f3)
# Override at run time:  docker run --user 1000:$(getent group dialout | cut -d: -f3)
ARG DIALOUT_GID=20
USER 1000:${DIALOUT_GID}

COPY --from=builder /serialcaster /serialcaster

ENTRYPOINT ["/serialcaster"]
