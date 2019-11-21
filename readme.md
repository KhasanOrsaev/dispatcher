Сервис-диспетчер rabbitmq->db
---
###Конфигурация
Файл *config.yaml.dist* заполнить и переименовать в **config.yaml**. Пример:
```yaml
time_out: 1    # периодичность выполнения скрипта в минутах,
services:
  source:
    host: localhost
    port: 5672
    user: user
    password: password
    queue:
      name: queue_name
      durable: true
      arguments:
        x-queue-mode: lazy
    workers: 2
  destination:
    host: db_host
    user: user
    password: password
    database: db
    queries:
      - query: insert into public.test(column1, column2, column3) values %s
        prepare_statement: (?,?,?)
    workers: 3
    bulk: 20
  crash:
    host: localhost
    port: 5672
    user: user
    password: password
    exchange:
      name: exchange_name
      routing_key: route_name
      type: topic
      durable: true
log:
  path: ./log.log
  format: "%s.%s message: %s context: %s extra: %s"
```