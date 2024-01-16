# ===== build stage ====
FROM golang:1.20.13-bullseye as builder

WORKDIR /app

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/gomod-cache \
    go mod download

COPY . .

RUN --mount=type=cache,target=/gomod-cache \
    --mount=type=cache,target=/go-cache \
    go build -trimpath -ldflags="-w -s" -o cmd/bin/main cmd/main.go

# ===== deploy stage ====
FROM golang:1.20.13-bullseye as deploy

WORKDIR /app

COPY --from=builder /app/cmd/bin/main .

COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.7.2 /lambda-adapter /opt/extensions/lambda-adapter

ENV PORT=3000
EXPOSE 3000

CMD ["/app/main"]
