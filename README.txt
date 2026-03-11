Веб-Клиент
Автор: Давлетов Динис
Использование: json-client [опции]

Ресурсы:
  posts     - 100 постов
  comments  - 500 комментариев
  albums    - 100 альбомов
  photos    - 5000 фото
  todos     - 200 задач
  users     - 10 пользователей

Опции:
  -resource string  Тип ресурса (по умолчанию \"posts\")
  -id int           Конкретный ID для загрузки (0 - все)
  -user int         Фильтр по ID пользователя
  -limit int        Лимит результатов (0 - без лимита)
  -workers int      Количество рабочих горутин (по умолчанию 3)
  -expand           Загружать связанные ресурсы (комментарии к постам и т.д.)
  -timeout duration Таймаут API (по умолчанию 10s)
  -help             Показать эту справку

Примеры:
  json-client -resource=posts -user=1 -limit=5
  json-client -resource=posts -id=3 -expand
  json-client -resource=users -expand
  json-client -resource=albums -user=1 -limit=3
  json-client -resource=todos -user=2 -limit=10

API: https://jsonplaceholder.typicode.com
