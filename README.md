

# Менеджер паролей GophKeeper

## Get started
- Сделать миграцию
```
goose -dir "./migrations" <db base driver> "<databaseDSN>" up
```
- Сборка
```
make build
```
- Создать и заполнить конфиги по шаблонам
```
config/example_client_cfg.yaml --> config/client_cfg.yaml
config/example_server_cfg.yaml --> config/server_cfg.yaml
```
- Генерация x.509 сертификатов и ключей.
```
make cert
```
- Запуск
```
Сервер
./cmd/server/server

CLI
./cmd/cli/cli
```