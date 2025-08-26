# netangels.dns-01
DNS challenge via NetAngels.ru for go-acme/lego

# .env
```.env
# .devcontainer/.env
NETANGELS_API_KEY=XXX
DEV_WATCH=1
```
# Время ожидания

Обновление DNS занимает
(на 22.08.2028):
48, 55, 64, 55, 60, 64, 54, 55, 67, 55

## Тесты

Часть тестов отключена
(`.skip`),
так как исполняются долго:
+ [Измерение времени обновления DNS](test/na.js)
+ [Запуск Lego](test/lego.js)

## See also
+ [NetAngels]
+ [go-acme/lego]
+ [Traefik]
+ [Let's Encrypt]


[go-acme/lego]: https://github.com/go-acme/lego
[Traefik]: https://traefik.io/traefik
[Let's Encrypt]: https://letsencrypt.org/
[NetAngels]: https://www.netangels.ru/
