FROM golang:1.25.2

RUN apt-get update && apt-get install -y make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --exclude=.env . .

RUN mkdir -p /app/migrations

RUN mkdir -p /app/logs
RUN touch /app/logs/app.log

COPY start.sh /start.sh
RUN chmod +x /start.sh

CMD ["/start.sh"]