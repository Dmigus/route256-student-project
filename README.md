# Курсовой проект 10-го потока route256 go middle


Требования, на основе которых реализовано приложение, размещены в директории `./docs`

Запуск всех сервисов и окружения в docker можно произвести командой `make run-all`
Остановку можно произвести командой `make stop-all`

UI окружения будут доступны по адресам:
- jaeger: localhost:16686
- kafka: localhost:8080
- prometheus: localhost:9090
- grafana: localhost:3000

Примеры запросов к сервисам cart и loms из Postman представлены в директории `./postman`
