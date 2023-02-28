# Тестовое задание whattofarm
## Условия
Создать Go приложение, состоящее из клиентской и серверной частей.

Приложение должно быть реализовано в виде CLI приложения на основе https://github.com/spf13/cobra

### Серверная часть

Должна предоставлять данные, предоставляемые публичным API биржи Binanсе, доступным по адресу https://api.binance.com/api/v3/ticker/price

Нужно реализовать два метода:

- GET

      $ curl http://localhost:3001/api/v1/rates?pairs=BTC-USDT,ETH-USDT
  
      { "ETH-USDT": 1780.123, "BTC-USDT": 46956.45 }
- POST

      $ curl http://localhost:3001/api/v1/rates?pairs=BTC-USDT,ETH-USDT
      { "ETH-USDT": 1780.123, "BTC-USDT": 46956.45 }

## Реализация
#### Запуск серверной части 

      $ go run ./server/main.go
      Listening at port 3001...
#### Клиентская часть
Возможные варианты:

- запуск без флага --pair выведет в консоль пары BTC-USDT и ETH-USDT

      $ go run ./client/main.go
      ETH-USDT:1639.10, BTC-USDT:23752.83

- запуск с флагом --pair=ETH-USDT выведет в консоль значение пары ETH-USDT

      $ go run ./client/main.go --pair=ETH-USDT
      1780.12

- запуск с флагом --pair==BTC-USDT выведет в консоль значение пары BTC-USDT

      $ go run ./client/main.go --pair=BTC-USDT
      24400.00

### Детали решения

1. При запуске сервера происходит инициализация сервера,роутера и переменных
отвечающих за значения пар.

2. При каждом обращении к серверу со стороны клиента происходит обновление значений переменных 
и выгрузка ответа в консоль.

3. Сервер завершает работу по сигналам Interrrupt или при ошибке.