# Karl'Go - Booking приложение на языке Go

<img width="1315" height="986" alt="Maps" src="https://github.com/user-attachments/assets/556390d7-e2c9-4128-8ef3-884077ec6ef6" />


## Проект создан на хакатоне "Моя профессия - ИТ"

За этот проект мы с командой на хакатоне "Моя профессия - ИТ" заняли 2-ое место. В команде я отвечал за бекенд составляющую (API,CRUD, интеграция с Yandex Geocoder, интеграция с LLM) , сервер, инфрастркутуру и настройку домена и SSL, swagger документацию. 

## 🛠 Технологический стек

### Основные технологии
![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)
![Beego](https://img.shields.io/badge/Beego-2.0+-1E1E20?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-4169E1?logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-24.0+-2496ED?logo=docker&logoColor=white)
![Docker Compose](https://img.shields.io/badge/Docker_Compose-2.23+-2496ED?logo=docker&logoColor=white)

### Frontend & DevOps
![Nginx](https://img.shields.io/badge/Nginx-1.25+-009639?logo=nginx&logoColor=white)
![React](https://img.shields.io/badge/React-18+-61DAFB?logo=react&logoColor=black)
![Vite](https://img.shields.io/badge/Vite-4.0+-646CFF?logo=vite&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-85EA2D?logo=swagger&logoColor=black)

### Интеграции
![Yandex Cloud](https://img.shields.io/badge/Yandex_Cloud-FFCC00?logo=yandex&logoColor=black)
![Let's Encrypt](https://img.shields.io/badge/Let's_Encrypt-003A70?logo=letsencrypt&logoColor=white)


## Идея и MVP

В рамках 2–3-дневного соревнования была выбрана задача, требующая быстрого взаимодействия клиента и сервера: генерация бизнес-описаний компаний и отображение их на карте, а так же бронирование кафе и ресторанов.

- За первый день был спроектирован REST API на Go с использованием Beego.
- Все эндпоинты были описаны в Swagger.
- Для хранения сущностей Company, Owner, User была подключена PostgreSQL (через client/orm).

## Backend

- **Beego + PostgreSQL:** реализованы модели организаций, владельцев и пользователей, операции CRUD, авторизация по JWT (разные ключи для владельцев и юзеров).
- **Интеграция OpenRouter и Яндекс.Геокодера:**
    - Был написан клиент для OpenRouter для генерации описаний компаний.
    - Асинхронный расчёт координат реализован через Yandex Geocoder, координаты сохранялись в базе (goroutine).
- **WebSocket-сервис:** для демонстрации real-time реализован простой эхо-сервер на Gorilla WebSocket.

## Docker + Nginx

- **Docker Compose:** в `docker-compose.yml` описаны контейнеры для backend (Go) и базы данных.
- **Nginx:** настроен обратный прокси, перенаправляющий `/api` на Go-сервис. Для HTTPS выполнена SSL-терминация (Let's Encrypt).

## Инфраструктура

- Виртуальный сервер: арендован и настроен VPS на Yandex Cloud (установлены Ubuntu, firewall (ufw), swap).
- Домен: зарегистрирован, настроена A-запись на IP VPS, автоматическое продление сертификата реализовано через certbot.
- **CI/CD (опционально):** реализован деплой по push в master через GitHub Actions или rsync+ssh с автоматическим перезапуском контейнера.

## Итог

За 2–3 дня было реализовано полностью работающее серверное приложение с защищённым API, генерацией контента и интеграциями с внешними сервисами.

Прокачаны навыки DevOps и Backend Go:  
От покупки домена и настройки сервера до написания backend-логики, интеграции внешних API, настройки контейнеров и развёртывания всего стека.

![Диплом финалиста](https://github.com/user-attachments/assets/4586c93c-6e54-4094-8a80-aaf392e10d00)

![Сертификат участника](https://github.com/user-attachments/assets/22779668-418a-4ff4-88c2-a8a1fb145ae5)

