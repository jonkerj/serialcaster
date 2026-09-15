# ── build stage ───────────────────────────────────────────────────────────────
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
ARG TARGETOS TARGETARCH

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /serialcaster .

# ── runtime stage ─────────────────────────────────────────────────────────────
FROM scratch

# Override at build time: docker build --build-arg DIALOUT_GID=$(getent group dialout | cut -d: -f3)
# Override at run time:  docker run --user 1000:$(getent group dialout | cut -d: -f3)
ARG DIALOUT_GID=20
USER 1000:${DIALOUT_GID}

COPY --from=builder /serialcaster /serialcaster

ENTRYPOINT ["/serialcaster"]
