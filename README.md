# Документация

## Как настроить и запустить приложение

### 1. настройка
Положите в корень проекта файл с названием `.env`.
Он **ОБЯЗАТЕЛЕН**, без него приложение откажется запускаться!

Переменные окружения:

| Название  | Тип               | Назначение                                                                                                   |
|-----------|-------------------|--------------------------------------------------------------------------------------------------------------|
| **GRPC_PORT** | `int` | **порт gRPC сервера**<br/>используется при запуске и при определении эндпоинта grpc-сервера при создании gateway |
| **HTTP_PORT** | `int` | **порт HTTP сервера**                                                                                            |

Хотите пример? Он есть в `./env.example`

### 2. запуск
```bash
make generate
make build
make run  # <-- можно обойтись одной командой
```

### 3. остановка

Приложение обрабатывает сигналы остановки и выполняет плавное завершение. В конце выводится "Server Stopped" с переводом строки