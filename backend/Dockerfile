FROM golang:1.23-alpine AS runner

WORKDIR /app

RUN apk add --no-cache curl git make

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/beego/bee/v2@latest


COPY . .



RUN go mod download

RUN go mod tidy

CMD ["bee", "run", "-gendoc=true", "-downdoc=true"]
