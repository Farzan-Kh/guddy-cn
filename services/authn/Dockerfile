FROM golang:1.24.6-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o authn .


FROM alpine:3.18

LABEL org.opencontainers.image.source=https://github.com/farzan-kh/guddy-cn

RUN apk --no-cache add ca-certificates

COPY --from=build /app/authn .

EXPOSE 8084

ENTRYPOINT ["./authn"]
