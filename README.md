# GOPH-KEEPER

## О проекте

Серверная часть **[./internal/server](./internal/server)**

Клиентская часть **[./internal/client](./internal/client)**

Взаимодействие осуществляется посредством **GRPC** протокола. Реализован TUI интерфейс для клиента. [Контракты GRPC](./contracts)

Серверная часть использует для хранения данных postgres, клиентская in-memory, без сторонних БД.

Для процедуры аутентификации используется JWT.

## Подготовленные бинарники

Можно скачать [тут](https://github.com/ktigay/goph-keeper/releases/latest)

Клиенты:
- [Linux клиент x64](https://github.com/ktigay/goph-keeper/releases/latest/download/goph-keeper)
- [Windows клиент x64](https://github.com/ktigay/goph-keeper/releases/latest/download/goph-keeper-windows-amd64.exe)
- [Mac клиент arm-x64](https://github.com/ktigay/goph-keeper/releases/latest/download/goph-keeper-darwin-arm64)
- [Mac клиент amd-x64](https://github.com/ktigay/goph-keeper/releases/latest/download/goph-keeper-darwin-amd64)

Сервер:
- [Linux сервер x64](https://github.com/ktigay/goph-keeper/releases/latest/download/goph-keeper-server)

Версию приложения можно получить командой
> ./goph-keeper -v

Справка по параметрам запуска
> ./goph-keeper -h

## Локальный запуск

### Билд докер-образов
> task build
 
### Билд бинарников
Сервер
> task build:server

Клиенты под все платформы

> task build:client:all

Бинарники можно найти тут ./bin 

### Запуск сервера и б.д.
> task up

### Запуск клиента из cli

> cd ./bin/client && ./goph-keeper
