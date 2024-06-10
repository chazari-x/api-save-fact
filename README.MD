Запуск Golang приложения `api-save-fact`
1. Перейдите в папку с проектом
2. Запустите приложение командой `go run main.go`
3. Для остановки приложения с ожиданием завершения всех запросов нажмите `Ctrl+C`
4. Для остановки приложения без ожидания завершения всех запросов нажмите `Ctrl+\`

Config:
1. В папке `etc` создайте файл `config.yaml`
2. Пример файла `config.yaml`
```
host: localhost
port: "8080"
bearer_token: kjergjl349u9032ut9r
href: https://development.kpi-drive.ru/_api/facts/save_fact
buffer: 2000
```

Endpoints:
1. POST /_api/facts/save_fact - сохранение факта. `body form/data`
```
period_start: 2024-05-01
period_end: 2024-05-31
period_key: month
indicator_to_mo_id: 227373
indicator_to_mo_fact_id: 0
value: 1
fact_time: 2024-05-31
is_plan: 0
auth_user_id: 40
comment: buffer Last_name
```
