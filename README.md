# Тестовое задание от EffectiveMobile Backend
# Сервис дополнения пользователей информацией

##### Автор: [Виноградов Данил](https://t.me/japsty) 

## Содержание
1. [Запуск](#запуск)
2. [Доступные методы](#доступные-методы)

### Запуск
#### Обычный запуск
```shell
  docker-compose up
```
#### Чистый запуск
```shell
  docker rm $(docker ps -a -q) && docker volume prune -f
  docker rmi -f em_test_task
  docker-compose up
```

### Доступные методы

*У проекта есть описание методов в [Postman](https://www.postman.com/japsty/workspace/danilvinogradov/collection/29141861-24c437d5-dd2e-4053-9767-36326e30a105?action=share&creator=32626453)*

#### **POST** /people
Метод добавления нового человека

Принимает такие параметры: name,surname,patronymic(опционально).
На основе переданных в него данных структура дополняется информацией из сторонних API, возраст и пол получаем по имени, национальность - по фамилии.

*Пример ниже показывает создание пользователя и конечную структуру*

*Принимаемая структура*
```json
{
    "name":"Danil",
    "surname":"Vinogradov",
    "patronymic":"Sergeevich"
}
```
*Дополненная структура*
```json
{
    "name": "Danil",
    "surname": "Vinogradov",
    "patronymic": "Sergeevich",
    "age": 57,
    "gender": "male",
    "nationality": "RU",
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z"
}
```

#### **GET** /people/"id"
Метод получения человека по его id в системе

Возвращает структуру, что представлена ниже

*Возвращаемая структура*
```json
{
    "id": 1,
    "name": "Danil",
    "surname": "Vinogradov",
    "patronymic": "Sergeevich",
    "age": 57,
    "gender": "male",
    "nationality": "RU",
    "created_at": "2024-01-29T08:45:58.188781Z",
    "updated_at": "2024-01-29T08:45:58.188781Z"
}
```

  
#### **GET** /people?page=1&per_page=5
Метод получения списка людей с пагинацией. Параметр page отвечает за страницу, per_page - за количество отображаемых на странице структур.

Возвращает набор структур, подобной той, что представлена ниже.

*Возвращаемая структура*
```json
{
    "id": 1,
    "name": "Danil",
    "surname": "Vinogradov",
    "patronymic": "Sergeevich",
    "age": 57,
    "gender": "male",
    "nationality": "RU",
    "created_at": "2024-01-29T08:45:58.188781Z",
    "updated_at": "2024-01-29T08:45:58.188781Z"
}
```

#### **GET** /people-filtered?page=1&per_page=5
Метод получения cписка людей с фильтрами и пагинацией

Возвращает структуры, соответствующие фильтру. Для фильтрации с обозначением границ "from" и "to" доступны поля "id" и "age", поля "nation" и "gender" доступны без границ, остальные элементы будут выводиться просто с пагинацией.

*Принимаемая структура*
```json
{
    "sort_by": "id",
    "from": 1,
    "to": 4
}
```

*Возвращаемая структура*
```json
{
    "id": 1,
    "name": "Danil",
    "surname": "Vinogradov",
    "patronymic": "Sergeevich",
    "age": 57,
    "gender": "male",
    "nationality": "RU",
    "created_at": "2024-01-29T08:45:58.188781Z",
    "updated_at": "2024-01-29T08:45:58.188781Z"
}
```

#### **PUT** /people/"id"
Метод обновления структуры по ее id.
"patronymic" - опционален, в случае его отсутствия это поле не будет выводиться.

Возвращает сообщение или ошибку.

*Принимаемая структура*
```json
{
    "name": "Vladimir",
    "surname":"Zabelin",
    "patronymic":"Vladimirovich"
}
```
*Возвращаемая структура*
```json
{
    "message": "Person updated successfully"
}
```

#### **DELETE** /people/"id"
Метод удаления структуры по id

Возвращает сообщение или ошибку.

*Возвращаемая структура*
```json
{
    "message": "Person was deleted"
}
```
