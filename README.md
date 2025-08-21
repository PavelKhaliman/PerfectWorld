# PerfectWorld (ГШ-созвездия)

- Проект сделан для учебной цели связанной с "прикруткой" go к php 

Небольшой Go‑сервис с двумя страницами:
- `/` — «Играть» (поле 3×3 сфер, инверсия соседей)
- `/solver.html` — «Решение» (поиск кратчайшей последовательности ходов)

CSS/оформление заимствовано из шаблона `blog`.

## Переменные окружения

- `PUBLIC_DIR` — путь к каталогу статики (по умолчанию `./public`)

## API

- `GET /api/state` → текущее поле (массив из 9 bool)
- `POST /api/toggle` (`index=1..9`) → инвертирует клетку и ортогональных соседей
- `POST /api/reset` → новое случайное поле
- `POST /api/solve` → для решения задач

## Реверс‑прокси (под /game/)

```nginx
location ^~ /game/ {
    proxy_pass         http://perfectworld:8080/; # или host.docker.internal:8080
    proxy_http_version 1.1;
    proxy_set_header   Host $host;
    proxy_set_header   X-Real-IP $remote_addr;
    proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header   X-Forwarded-Proto $scheme;
}
```

В проекте пути в HTML/JS относительные (`css/`, `js/`, `assets/`, `api/...`), чтобы корректно работать под префиксом `/game/`.

## Docker Compose (вариант с блогом)

```yaml
services:
  perfectworld:
    image: alpine:3.20
    container_name: perfectworld
    restart: unless-stopped
    working_dir: /app
    volumes:
      - /root/perfectworld:/app:ro
    command: ["/app/perfectworld"]
    expose:
      - "8080"
    networks:
      - laravel_network
```




