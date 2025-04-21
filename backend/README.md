Для использования этого Dockerfile:

Поместите его в корень вашего Go-проекта

Соберите образ: 
```bash
docker build -t grafana-deploy .
```


Запустите: 
```bash
docker run -p 8080:8080 grafana-deploy
```