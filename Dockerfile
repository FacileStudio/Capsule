FROM oven/bun:1 AS client-build
WORKDIR /client
COPY apps/client/package.json apps/client/bun.lock* ./
RUN bun install --frozen-lockfile
COPY apps/client/ .
RUN bun run build

# 1.26 matches go.mod. The builder was on 1.25 while go.mod said 1.26, which
# made the toolchain silently download a newer one mid-build.
FROM golang:1.26-alpine AS api-build

ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /repo/apps/api

COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download

COPY apps/api ./

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o bin/api .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=api-build /repo/apps/api/bin/api /api
COPY --from=client-build /client/build /client

# The :nonroot base sets WorkingDir=/home/nonroot, so a relative ./client would
# resolve there and the SPA would silently not be served. Be explicit.
ENV CLIENT_DIR=/client

EXPOSE 4000

USER nonroot:nonroot

ENTRYPOINT ["/api"]
