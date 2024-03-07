FROM golang:1.22.0 as builder
WORKDIR /builder

RUN mkdir logs

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/main.go

FROM scratch as prod
WORKDIR /prod

COPY config/config.yaml ./config/
COPY logs ./logs
COPY --from=builder /builder/main .

EXPOSE 9999
CMD [ "./main" ]