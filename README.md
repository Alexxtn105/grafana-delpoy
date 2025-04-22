## Видео
https://youtu.be/LyocQr7cN-0

## Поднять образ (локально)
```bash
docker compose up
```

## Поднять образ (прод) - с параметром `-d` - чтобы все работало в фоновом режиме
```bash
docker compose up --build -d
```

## Посмотреть метрики Prometheus (в браузере)
http://localhost:9090/targets?search=


## Посмотреть логи в Promtail (в браузере)
http://localhost:9080/targets?search=

## Посмотреть Grafana
http://localhost:3000
Там же можно проверить наличие источников данных prometheus и Loki:
http://localhost:3000/connections/datasources

## Посмотреть метрики
http://localhost:8080/metrics

## Нагрузочное тестирование с ab:

```bash
ab -k -c 5 -n 20000 'http://localhost:8080/' & \
ab -k -c 5 -n 2000 'http://localhost:8080/status/400' & \
ab -k -c 5 -n 3000 'http://localhost:8080/status/409' & \
ab -k -c 5 -n 5000 'http://localhost:8080/status/500' & \
ab -k -c 50 -n 5000 'http://localhost:8080/status/200?seconds_sleep=1' & \
ab -k -c 50 -n 2000 'http://localhost:8080/status/200?seconds_sleep=2'
```

Или одной строкой:
```bash
ab -k -c 5 -n 20000 'http://localhost:8080/' & ab -k -c 5 -n 2000 'http://localhost:8080/status/400' & ab -k -c 5 -n 3000 'http://localhost:8080/status/409' & ab -k -c 5 -n 5000 'http://localhost:8080/status/500' & ab -k -c 50 -n 5000 'http://localhost:8080/status/200?seconds_sleep=1' & ab -k -c 50 -n 2000 'http://localhost:8080/status/200?seconds_sleep=2'
```

## Пример дашборда
Находится в папке /grafana/example-dashboard.json