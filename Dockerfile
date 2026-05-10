FROM golang:1.25 AS build-stage

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /binary cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /binary /binary
COPY ./static /static

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/binary"]
