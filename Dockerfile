FROM golang:1.20 AS build

WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/*.go ./
COPY backend/web ./web

RUN CGO_ENABLED=0 go build -o /sports-app .

FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=build /sports-app /sports-app
COPY backend/web /web

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/sports-app"]