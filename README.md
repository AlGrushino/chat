# chat

git clone <repository-url>
cd chat

миграции выполняются автоматически
docker-compose up --build

приложение будет доступно по адресу:
http://localhost:8080

ручное выполнение миграций:
make migrate или go run cmd/main.go --migrate

запустить тесты:
make test
