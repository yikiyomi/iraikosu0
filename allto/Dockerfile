ARG GO_VERSION=1.26.1

# 阶段1：下载依赖
FROM golang:${GO_VERSION} AS deps
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

# 阶段2：编译
FROM deps AS build
COPY . ./
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/allto .

# 阶段3：运行
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -H -s /sbin/nologin app
WORKDIR /app
COPY --from=build /out/allto /app/allto
COPY --from=build /src/config.docker.yaml /app/config.yaml
COPY --from=build /src/frontend /app/frontend
USER app
EXPOSE 8080
ENTRYPOINT ["/app/allto"]