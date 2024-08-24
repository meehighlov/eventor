# grats

## Запуск

1. создаем .env файл и добавляем его в /cmd (пример .env.example)
2. запускаем миграции
   - устанавливаем goose, например: brew install goose
   - из каталога /cmd запускаем команду
   ```shell
   goose -dir=../migrations sqlite3 eventor.db up
   ```
3. запускаем бота из каталога /cmd:
   ```shell
   go run main.go
   ```
