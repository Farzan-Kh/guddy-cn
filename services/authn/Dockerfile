FROM golang:1.24.6-alpine AS build
WORKDIR /src
COPY . .
RUN apk add --no-cache gcc musl-dev
RUN go build -o /authn ./

FROM alpine:3.18
LABEL org.opencontainers.image.source=https://github.com/farzan-kh/guddy-cn
COPY --from=build /authn /authn
EXPOSE 8084
ENTRYPOINT ["/authn"]
