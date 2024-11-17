FROM golang:1.23-alpine AS builder

WORKDIR /app

# cache deps
COPY go.mod .
COPY go.sum .
RUN go mod download

# build
COPY . .
RUN CGO_ENABLED=0 go build

FROM golang:1.23-alpine AS runner
WORKDIR /app
COPY --from=builder /app/simplebank .
ENV PATH="/app:${PATH}"
CMD ["simplebank"]
