# fast_mtproxy
Быстрое развёртывание официального MTProxy сервера с поддержкой Fake TLS

## Установка
### Загрузка исполняемого файла
```bash
curl -L -o install https://git.io/JeOSr && chmod +x install
```

### Быстрая установка
```bash
./install
```

По умолчанию: порт 443, генерируется рандомный секретный ключ, домен www.google.com.

### Пример
```bash
./install -p=8888 -H=443 -S=25c8dfee81acdadaff3a338a10db8497 -P=<получить тут @MTProxybot> -D=www.google.com
```
<code>-p</code> - локальный порт для статистики (необязательно);<br>
<code>-H</code> - порт для подключения;<br>
<code>-S</code> - секретный ключ;<br>
<code>-P</code> - рекламный тег;<br>
<code>-D</code> - домен с поддержкой TLS 1.3 (проверить можно <a href="https://www.cdn77.com/tls-test">тут</a>).
