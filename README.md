<h1> Web programming and mobile applications</h1>
<details>
  <summary><h2>🎄 Lab 1 - new-year-counter</h2></summary>
  
  <p>Использовали:</p>
  <ul>
    <li>Go</li>
    <li>Docker</li>
  </ul>
  
  <p>Инструкция</p>
  <ul>
    <li>Скачать</li>
    <li>Собрать проект командой - <b>docker build -t new-year-counter .</b></li>
    <li>Запустить проект - <b>docker run -d -p 3000:3000 --name new-year-app new-year-counter</b></li>
    <li>Запрос на получение результата - <b>curl http://localhost:3000</b></li>
  </ul>
  
  <p>Быстрый обзор результата</p>
  <h3>Что отображается в docker</h3>
  <img src="./assets/lab1/img1.png" alt="Что появляется в Docker после build and run" width="600"/>
  <br>
  <h3>Что отображается в браузере (если всё хорошо)</h3>
  <img src="./assets/lab1/img2.png" alt="Что отображается в браузере localhost:3000" width="600"/>
  
  <i>Приносим изменения за отсутствие комментариев в коде</i>
</details>

<h1>Rugram API - Social Media Backend Service</h1>

<details>
  <summary><h2>📱 Lab 2 - rugram-api (Base API)</h2></summary>
  
  <h3>Используемые технологии:</h3>
  <ul>
    <li><b>Go 1.21+</b> - основной язык разработки</li>
    <li><b>Gin Framework</b> - HTTP веб-фреймворк</li>
    <li><b>PostgreSQL 16</b> - реляционная база данных</li>
    <li><b>Docker & Docker Compose</b> - контейнеризация и оркестрация</li>
    <li><b>lib/pq</b> - драйвер PostgreSQL для Go</li>
    <li><b>godotenv</b> - управление переменными окружения</li>
    <li><b>google/uuid</b> - генерация уникальных идентификаторов</li>
  </ul>
  
  <h3>Архитектура проекта</h3>
  <ul>
    <li><b>Models</b> - данные</li>
    <li><b>Repository</b> - доступ к данным</li>
    <li><b>Service</b> - бизнес-логики</li>
    <li><b>Handlers</b> - HTTP обработчики</li>
    <li><b>DTO</b> - объекты передачи данных</li>
    <li><b>Config</b> - управление конфигурацией через .env</li>
  </ul>
  
  <h3>REST API Endpoints (Базовая версия)</h3>
  <table>
    <tr>
      <th>Method</th>
      <th>Endpoint</th>
      <th>Description</th>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/posts</td>
      <td>Получить все посты (с пагинацией)</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/posts/:id</td>
      <td>Получить пост по ID</td>
    </tr>
    <tr>
      <td>POST</td>
      <td>/api/v1/posts</td>
      <td>Создать новый пост</td>
    </tr>
    <tr>
      <td>PUT</td>
      <td>/api/v1/posts/:id</td>
      <td>Полностью обновить пост</td>
    </tr>
    <tr>
      <td>PATCH</td>
      <td>/api/v1/posts/:id</td>
      <td>Частично обновить пост</td>
    </tr>
    <tr>
      <td>DELETE</td>
      <td>/api/v1/posts/:id</td>
      <td>Мягкое удаление поста</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/posts/user/:userId</td>
      <td>Получить посты пользователя</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/health</td>
      <td>Health check сервиса</td>
    </tr>
   </table>
</details>

<details>
  <summary><h2>Lab 3 - rugram-api (Extended with Auth & Users)</h2></summary>
  
  <h3>Новый функционал</h3>
  <ul>
    <li><b>Аутентификация</b> - JWT токены (Access + Refresh)</li>
    <li><b>Управление пользователями</b> - регистрация, профиль, обновление данных</li>
    <li><b>OAuth 2.0</b> - Вход через Яндекс и ВКонтакте</li>
    <li><b>Безопасность</b> - bcrypt для паролей, хеширование токенов</li>
    <li><b>Сессии</b> - управление несколькими сессиями, logout-all</li>
    <li><b>Soft Delete</b> - безопасное удаление пользователей</li>
  </ul>
  
  <h4>Аутентификация и пользователи</h4>
  <table>
    <tr>
      <th>Method</th>
      <th>Endpoint</th>
      <th>Description</th>
      <th>Auth Required</th>
    </tr>
    <tr>
      <td>POST</td>
      <td>/api/v1/auth/register</td>
      <td>Регистрация нового пользователя</td>
      <td>❌</td>
    </tr>
    <tr>
      <td>POST</td>
      <td>/api/v1/auth/login</td>
      <td>Вход в систему (устанавливает cookies)</td>
      <td>❌</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/auth/whoami</td>
      <td>Получить информацию о текущем пользователе</td>
      <td>✅</td>
    </tr>
    <tr>
      <td>POST</td>
      <td>/api/v1/auth/refresh</td>
      <td>Обновить access token через refresh token</td>
      <td>❌</td>
    </tr>
    <tr>
      <td>POST</td>
      <td>/api/v1/auth/logout</td>
      <td>Выход из текущей сессии</td>
      <td>✅</td>
    </tr>
    <tr>
      <td>POST</td>
      <td>/api/v1/auth/logout-all</td>
      <td>Выход из всех сессий пользователя</td>
      <td>✅</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/users/:id</td>
      <td>Получить пользователя по ID</td>
      <td>✅</td>
    </tr>
    <tr>
      <td>PUT/PATCH</td>
      <td>/api/v1/users/:id</td>
      <td>Обновить данные пользователя</td>
      <td>✅</td>
    </tr>
    <tr>
      <td>DELETE</td>
      <td>/api/v1/users/:id</td>
      <td>Мягкое удаление аккаунта</td>
      <td>✅</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/users</td>
      <td>Список пользователей (с пагинацией)</td>
      <td>✅</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/users/email/:email</td>
      <td>Поиск по email</td>
      <td>✅</td>
    </tr>
  </table>
  
  <h4>OAuth 2.0 провайдеры</h4>
  <table>
    <tr>
      <th>Method</th>
      <th>Endpoint</th>
      <th>Description</th>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/auth/oauth/yandex</td>
      <td>Вход через Яндекс (редирект)</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/auth/oauth/vk</td>
      <td>Вход через ВКонтакте (редирект)</td>
    </tr>
  </table>
  
  <h4>Посты (расширенные)</h4>
  <table>
    <tr>
      <th>Method</th>
      <th>Endpoint</th>
      <th>Description</th>
      <th>Auth Required</th>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/posts</td>
      <td>Все посты с пагинацией</td>
      <td>error</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/posts/:id</td>
      <td>Пост по ID</td>
      <td>error</td>
    </tr>
    <tr>
      <td>POST</td>
      <td>/api/v1/posts</td>
      <td>Создать пост (привязан к user_id)</td>
      <td>ok</td>
    </tr>
    <tr>
      <td>PUT/PATCH</td>
      <td>/api/v1/posts/:id</td>
      <td>Обновить свой пост</td>
      <td>ok</td>
    </tr>
    <tr>
      <td>DELETE</td>
      <td>/api/v1/posts/:id</td>
      <td>Удалить свой пост</td>
      <td>ok</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/posts/user/:userId</td>
      <td>Посты пользователя</td>
      <td>error</td>
    </tr>
  </table>
  
  <h3>Примеры запросов</h3>
  
  <h4>Регистрация</h4>
  <pre><code>curl -X POST http://localhost:4200/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "phone": "+79991234567"
  }'</code></pre>
  
  <h4>Вход (устанавливает cookies)</h4>
  <pre><code>curl -X POST http://localhost:4200/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }' \
  -c cookies.txt</code></pre>
  
  <h4>Кто я (текущий пользователь)</h4>
  <pre><code>curl http://localhost:4200/api/v1/auth/whoami \
  -b cookies.txt</code></pre>
  
  <h4>Создать пост (авторизованный)</h4>
  <pre><code>curl -X POST http://localhost:4200/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "your-uuid-here",
    "title": "Мой первый пост!",
    "description": "Создано через API",
    "status": "active"
  }' \
  -b cookies.txt</code></pre>
  
  <h4>Обновить профиль</h4>
  <pre><code>curl -X PATCH http://localhost:4200/api/v1/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com",
    "phone": "+79876543210"
  }' \
  -b cookies.txt</code></pre>
  
  <h4>Выход из всех сессий</h4>
  <pre><code>curl -X POST http://localhost:4200/api/v1/auth/logout-all \
  -b cookies.txt</code></pre>
  
  <h3>Модели данных</h3>
  
  <h4>Users Table</h4>
  <pre><code>{
  "id": "uuid",
  "email": "user@example.com",
  "phone": "+79991234567",
  "yandex_id": "optional",
  "vk_id": "optional",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "deleted_at": "timestamp (soft delete)"
}</code></pre>
  
  <h4>User Tokens Table</h4>
  <pre><code>{
  "id": "uuid",
  "user_id": "uuid",
  "token_hash": "sha256 hash",
  "token_type": "access | refresh",
  "expires_at": "timestamp",
  "revoked": "boolean"
}</code></pre>

  <h3>🔧 Environment Variables</h3>
  
  <pre><code># Database
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=rugram_user
    DB_PASSWORD=rugram_password
    DB_NAME=rugram_db

    # App
    APP_PORT=4200
    APP_ENV=development

    # Pagination
    DEFAULT_PAGE=1
    DEFAULT_LIMIT=10
    MAX_LIMIT=100

    # JWT Secrets (измените в production!)
    JWT_ACCESS_SECRET=your-super-secret-access-key-here
    JWT_REFRESH_SECRET=your-super-secret-refresh-key-here

    # OAuth Yandex
    YANDEX_CLIENT_ID=your_yandex_client_id
    YANDEX_CLIENT_SECRET=your_yandex_client_secret
    YANDEX_REDIRECT_URI=http://localhost:4200/api/v1/auth/oauth/yandex/callback

    # OAuth VK
    VK_CLIENT_ID=your_vk_client_id
    VK_CLIENT_SECRET=your_vk_client_secret
    VK_REDIRECT_URI=http://localhost:4200/api/v1/auth/oauth/vk/callback

</code></pre>

  <h3>Инструкция по запуску</h3>
  
  <p><b>1. Клонировать и настроить</b></p>
  <pre><code>git clone https://github.com/yourusername/rugram-api.git
cd rugram-api
cp .env.example .env
# Отредактируйте .env, добавьте JWT_SECRET</code></pre>
  
  <p><b>2. Запуск с Docker Compose</b></p>
  <pre><code># Собрать и запустить
docker-compose up -d --build

# Проверить логи

docker-compose logs -f api

# Выполнить миграции (автоматически при старте)</code></pre>

  <p><b>3. Проверка работы</b></p>
  <pre><code># Health check
curl http://localhost:4200/health

# Регистрация

curl -X POST http://localhost:4200/api/v1/auth/register \
 -H "Content-Type: application/json" \
 -d '{"email":"test@test.com","password":"test123"}'

# Логин

curl -X POST http://localhost:4200/api/v1/auth/login \
 -H "Content-Type: application/json" \
 -d '{"email":"test@test.com","password":"test123"}' \
 -c cookies.txt

# Проверка whoami

curl http://localhost:4200/api/v1/auth/whoami -b cookies.txt</code></pre>

  <h3>Скриншоты работы (Lab 3)</h3>
  
  <h4>Регистрация нового пользователя</h4>
  <img src="./assets/lab3/register.png" alt="POST /auth/register создание пользователя" width="600"/>
  
  <h4>Вход и установка cookies</h4>
  <img src="./assets/lab3/login.png" alt="POST /auth/login с сохранением cookies" width="600"/>
  
  <h4>Получение информации о текущем пользователе</h4>
  <img src="./assets/lab3/whoami.png" alt="GET /auth/whoami возвращает данные пользователя" width="600"/>
  
  <h4>Создание поста авторизованным пользователем</h4>
  <img src="./assets/lab3/create-post-auth.png" alt="POST /posts с авторизацией" width="600"/>
  
  <h4>Обновление профиля пользователя</h4>
  <img src="./assets/lab3/update-user.png" alt="PATCH /users/:id обновление данных" width="600"/>
  
  <h4>Список всех пользователей (пагинация)</h4>
  <img src="./assets/lab3/users-list.png" alt="GET /users с пагинацией" width="600"/>
  
  <h4>Docker контейнеры</h4>
  <img src="./assets/lab3/docker-containers.png" alt="Запущенные контейнеры" width="600"/>
  
  <h4>Таблицы в PostgreSQL</h4>
  <img src="./assets/lab3/postgres-tables.png" alt="Структура БД: users, user_tokens, posts" width="600"/>
  
  <h4>Логи аутентификации</h4>
  <img src="./assets/lab3/auth-logs.png" alt="Логи с JWT и OAuth операциями" width="600"/>
  
  <h3>Решение проблем</h3>
  
  <p><b>Ошибка: "relation already exists" при миграциях</b></p>
  <ul>
    <li>Используйте <code>CREATE TABLE IF NOT EXISTS</code> и <code>CREATE INDEX IF NOT EXISTS</code></li>
    <li>Очистите volume: <code>docker-compose down -v</code></li>
    <li>Удалите таблицу <code>schema_migrations</code> если используется</li>
  </ul>
  
  <p><b>Ошибка: "invalid or expired token"</b></p>
  <ul>
    <li>Проверьте системное время на сервере</li>
    <li>Убедитесь что JWT_SECRET одинаковый при создании и проверке</li>
    <li>Access token живет 15 минут, используйте /refresh</li>
  </ul>
  
  <p><b>OAuth не работает</b></p>
  <ul>
    <li>Зарегистрируйте приложение в Яндекс.OAuth и VK API</li>
    <li>Укажите правильные Redirect URIs</li>
    <li>Проверьте переменные окружения YANDEX_CLIENT_ID и др.</li>
  </ul>

</details>

<details>
  <summary><h2>Lab 5 - rugram-api (Redis Cache & Session Management)</h2></summary>
  
  <h3>Новый функционал</h3>
  <ul>
    <li><b>Redis Cache</b> - кеширование часто запрашиваемых данных</li>
    <li><b>Умная инвалидация кеша</b> - автоматическое удаление при создании/обновлении/удалении</li>
    <li><b>JTI (JWT ID) в Redis</b> - мгновенный отзыв access токенов при logout</li>
    <li><b>Кеширование профилей пользователей</b> - снижение нагрузки на БД</li>
    <li><b>Кеширование списков постов</b> - с учетом пагинации</li>
    <li><b>TTL управление</b> - автоматическое удаление устаревших данных</li>
    <li><b>Безопасное хранение</b> - только JTI, без паролей и чувствительных данных</li>
  </ul>
  
  <h4>Структура ключей Redis</h4>
  <pre><code># Посты
rugram:posts:list:page:1:limit:10
rugram:posts:user:{userId}:list:page:1:limit:10
rugram:posts:item:{postId}

# Пользователи

rugram:users:profile:{userId}
rugram:users:email:{email}
rugram:users:list:page:1:limit:10

# Сессии (JTI токенов)

rugram:auth:user:{userId}:access:{jti}
rugram:auth:user:{userId}:refresh:{jti}</code></pre>

  <h4>Технологии</h4>
  <ul>
    <li><b>Redis 7 Alpine</b> - In-memory data store</li>
    <li><b>go-redis/v9</b> - Redis клиент для Go</li>
    <li><b>JWT with JTI</b> - Уникальные идентификаторы токенов</li>
    <li><b>Cache-Aside стратегия</b> - Lazy loading кеша</li>
  </ul>
  
  <h3>Примеры запросов с кешированием</h3>
  
  <h4>Проверка кеширования постов</h4>
  <pre><code># Первый запрос - загрузит из БД и сохранит в кеш
time curl -X GET "http://localhost:4200/api/v1/posts?page=1&limit=10" -b cookies.txt

# Второй запрос - возьмет из Redis (должен быть быстрее)

time curl -X GET "http://localhost:4200/api/v1/posts?page=1&limit=10" -b cookies.txt</code></pre>

  <h4>Создание поста (инвалидация кеша)</h4>
  <pre><code># После этого запроса кеш списков будет очищен
curl -X POST http://localhost:4200/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "'"$USER_ID"'",
    "title": "Новый пост",
    "description": "Этот пост очистит кеш",
    "status": "active"
  }' \
  -b cookies.txt</code></pre>
  
  <h4>Проверка отзыва токена через Redis</h4>
  <pre><code># Логин и сохранение cookies
curl -X POST http://localhost:4200/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"SecurePass123"}' \
  -c cookies.txt

# Проверка whoami - работает

curl http://localhost:4200/api/v1/auth/whoami -b cookies.txt

# Выход - удаляет JTI из Redis

curl -X POST http://localhost:4200/api/v1/auth/logout -b cookies.txt

# Повторный запрос с тем же токеном - 401 Unauthorized

curl http://localhost:4200/api/v1/auth/whoami -b cookies.txt</code></pre>

  <h3>Redis CLI команды для проверки</h3>
  
  <pre><code># Подключение к Redis
docker exec -it rugram_redis redis-cli -a redis_secure_password_change_in_prod

# Просмотр всех ключей кеша

KEYS rugram:\*

# Просмотр содержимого кеша постов

GET "rugram:posts:list:page:1:limit:10"

# Проверка TTL (Time To Live)

TTL "rugram:posts:list:page:1:limit:10"

# Просмотр активных сессий (JTI токенов)

KEYS "rugram:auth:user:_:access:_"

# Ручная инвалидация кеша (для тестов)

DEL "rugram:posts:list:page:1:limit:10"

# Удаление всех ключей по паттерну

redis-cli -a your_password KEYS "rugram:posts:\*" | xargs redis-cli -a your_password DEL

# Мониторинг операций в реальном времени

MONITOR</code></pre>

  <h3>Скриншоты работы (Lab 5)</h3>
  
  <h4>Redis контейнер в Docker</h4>
  <img src="./assets/lab5/redis-container.png" alt="Redis контейнер запущен" width="600"/>
  
  <h4>Ключи кеша в Redis</h4>
  <img src="./assets/lab5/redis-keys.png" alt="Просмотр ключей через redis-cli" width="600"/>
  
  <h4>Кеширование списка постов</h4>
  <img src="./assets/lab5/posts-cache.png" alt="Кешированный ответ постов" width="600"/>
  
  <h4>JTI токены в Redis (активные сессии)</h4>
  <img src="./assets/lab5/jti-tokens.png" alt="JTI токены пользователя в Redis" width="600"/>
  
  <h4>Сравнение времени ответа (с/без кеша)</h4>
  <img src="./assets/lab5/performance-comparison.png" alt="Сравнение производительности" width="600"/>
  
  <h4>Логи приложения с кешированием</h4>
  <img src="./assets/lab5/cache-logs.png" alt="Логи показывают Cache HIT/MISS" width="600"/>
  
  <h3>Проверка работы кеша</h3>
  
  <p><b>1. Запустите приложение</b></p>
  <pre><code>docker-compose up -d --build</code></pre>
  
  <p><b>2. Создайте тестовые данные</b></p>
  <pre><code># Регистрация
curl -X POST http://localhost:4200/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"cache@test.com","password":"test123"}'

# Логин

curl -X POST http://localhost:4200/api/v1/auth/login \
 -H "Content-Type: application/json" \
 -d '{"email":"cache@test.com","password":"test123"}' \
 -c cookies.txt

# Создайте несколько постов

for i in {1..5}; do
curl -X POST http://localhost:4200/api/v1/posts \
 -H "Content-Type: application/json" \
 -d "{\"user_id\":\"$USER_ID\",\"title\":\"Post $i\",\"status\":\"active\"}" \
 -b cookies.txt
done</code></pre>

  <p><b>3. Проверьте Redis кеш</b></p>
  <pre><code>docker exec -it rugram_redis redis-cli -a redis_secure_password_change_in_prod
KEYS rugram:*
GET "rugram:posts:list:page:1:limit:10"</code></pre>
  
  <h3>Полезные команды</h3>
  
  <pre><code>
  # Посмотреть статистику Redis
  docker exec -it rugram_redis redis-cli -a your_password INFO stats

# Мониторинг запросов в реальном времени

docker exec -it rugram_redis redis-cli -a your_password MONITOR

# Очистить все кеши (только для тестов)

docker exec -it rugram_redis redis-cli -a your_password FLUSHDB

# Бэкап Redis данных

docker exec rugram_redis redis-cli -a your_password SAVE

# Просмотр логов приложения с фильтром по кешу

docker logs rugram_api 2>&1 | grep -i cache
</code></pre>

</details>

# Лабораторная работа №6 - Миграция на MongoDB

## Выполнил: [Ваше имя]
## Группа: [Ваша группа]

---

<details open>
  <summary><h2>Lab 6 - rugram-api (MongoDB Migration)</h2></summary>
  
  <h3>Новый функционал</h3>
  <ul>
    <li><b>Миграция с PostgreSQL на MongoDB</b> - полная замена реляционной БД на документоориентированную</li>
    <li><b>Документная модель данных</b> - гибкие схемы для пользователей, постов и токенов</li>
    <li><b>ObjectID вместо UUID</b> - нативные MongoDB идентификаторы</li>
    <li><b>Встроенные индексы</b> - оптимизация запросов на уровне коллекций</li>
    <li><b>Mongo Express</b> - веб-админка для управления БД</li>
    <li><b>Сохранение функциональности</b> - все API эндпоинты работают как прежде</li>
    <li><b>Soft Delete через поле deletedAt</b> - механизм мягкого удаления сохранен</li>
  </ul>
  
  <h4>Сравнение моделей данных</h4>
  
  <p><b>PostgreSQL (было):</b></p>
  <pre><code>CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    created_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    title VARCHAR(200),
    created_at TIMESTAMP
);</code></pre>

  <p><b>MongoDB (стало):</b></p>
  <pre><code>// Коллекция users
{
    "_id": ObjectId("67f4b8a9d2f4a12c34567890"),
    "email": "user@example.com",
    "password_hash": "bcrypt_hash",
    "phone": "+79123456789",
    "yandex_id": "12345",
    "created_at": ISODate("2024-01-01T00:00:00Z"),
    "updated_at": ISODate("2024-01-01T00:00:00Z"),
    "deleted_at": null
}

// Коллекция posts
{
    "_id": ObjectId("67f4b8a9d2f4a12c34567891"),
    "user_id": "67f4b8a9d2f4a12c34567890",  // Ссылка на _id пользователя
    "title": "Мой первый пост",
    "description": "Текст поста",
    "status": "active",
    "likes_count": 0,
    "created_at": ISODate("2024-01-01T00:00:00Z"),
    "deleted_at": null
}</code></pre>

  <h4>Технологии</h4>
  <ul>
    <li><b>MongoDB 6</b> - Документоориентированная СУБД</li>
    <li><b>MongoDB Go Driver 1.13.1</b> - Официальный драйвер для Go</li>
    <li><b>Mongo Express</b> - Веб-интерфейс для администрирования</li>
    <li><b>BSON</b> - Бинарный формат хранения данных</li>
  </ul>
  
  <h3>Преимущества MongoDB в проекте</h3>
  
  <h4>1. Гибкая схема данных</h4>
  <pre><code>// OAuth пользователи не нуждаются в password_hash
{
    "_id": ObjectId("..."),
    "email": "oauth@yandex.ru",
    "yandex_id": "123456",
    // password_hash отсутствует - это допустимо!
    "created_at": ISODate("...")
}</code></pre>

  <h4>2. Встроенные массивы (будущие возможности)</h4>
  <pre><code>// Можно хранить комментарии прямо в посте
{
    "_id": ObjectId("..."),
    "title": "Пост с комментариями",
    "comments": [
        {"user_id": "...", "text": "Отличный пост!", "created_at": ISODate("...")},
        {"user_id": "...", "text": "Согласен!", "created_at": ISODate("...")}
    ]
}</code></pre>

  <h4>3. Атомарные обновления</h4>
  <pre><code>// Инкремент лайков без дополнительных запросов
db.posts.updateOne(
    {"_id": ObjectId("...")},
    {"$inc": {"likes_count": 1}}
)</code></pre>

  <h4>Индексы MongoDB</h4>
  <pre><code>// Уникальный индекс на email
db.users.createIndex({"email": 1}, {unique: true})

// Составной индекс для фильтрации
db.posts.createIndex({"user_id": 1, "status": 1, "deleted_at": 1})

// TTL индекс для автоматической очистки токенов
db.user_tokens.createIndex({"expires_at": 1}, {expireAfterSeconds: 0})</code></pre>

  <h3>Примеры запросов к MongoDB</h3>
  
  <h4>Создание пользователя</h4>
  <pre><code>curl -X POST http://localhost:4200/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mongodb@test.com",
    "password": "SecurePass123",
    "phone": "+79123456789"
  }'

# В MongoDB Express можно увидеть:
# {
#   "_id": ObjectId("67f4b8a9d2f4a12c34567890"),
#   "email": "mongodb@test.com",
#   "phone": "+79123456789",
#   "password_hash": "$2a$10$...",
#   "created_at": ISODate("2024-01-01T00:00:00Z")
# }</code></pre>

  <h4>Создание поста пользователем</h4>
  <pre><code>curl -X POST http://localhost:4200/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "67f4b8a9d2f4a12c34567890",
    "title": "Тестовый пост в MongoDB",
    "description": "Проверка работы с документной БД",
    "status": "active"
  }' \
  -b cookies.txt</code></pre>

  <h4>Получение постов с пагинацией</h4>
  <pre><code>curl "http://localhost:4200/api/v1/posts?page=1&limit=10" -b cookies.txt

# MongoDB запрос:
# db.posts.find({"deleted_at": null})
#   .sort({"created_at": -1})
#   .skip(0)
#   .limit(10)</code></pre>

  <h3>Mongo Express - Админ панель</h3>
  
  <h4>Доступ к админке</h4>
  <pre><code>URL: http://localhost:8081
Login: admin
Password: admin_password</code></pre>

  <h4>Скриншоты Mongo Express</h4>
  
  <p><b>Главная страница со списком БД</b></p>
  <img src="./assets/lab6/mongo-express-dashboard.png" alt="Mongo Express Dashboard" width="800"/>
  
  <p><b>Коллекция users</b></p>
  <img src="./assets/lab6/users-collection.png" alt="Users collection в MongoDB" width="800"/>
  
  <p><b>Просмотр документа пользователя</b></p>
  <img src="./assets/lab6/user-document.png" alt="Документ пользователя" width="800"/>
  
  <p><b>Коллекция posts</b></p>
  <img src="./assets/lab6/posts-collection.png" alt="Posts collection" width="800"/>
  
  <p><b>Индексы коллекции posts</b></p>
  <img src="./assets/lab6/posts-indexes.png" alt="Индексы MongoDB" width="800"/>

  <h3>MongoDB CLI команды</h3>
  
  <pre><code># Подключение к MongoDB Shell
docker exec -it rugram_mongo mongosh -u rugram_user -p rugram_password --authenticationDatabase admin

# Просмотр всех БД
show dbs

# Выбор БД
use rugram_db

# Просмотр коллекций
show collections

# Поиск пользователей
db.users.find().pretty()

# Поиск по email
db.users.find({"email": "mongodb@test.com"}).pretty()

# Поиск активных постов пользователя
db.posts.find({
    "user_id": "67f4b8a9d2f4a12c34567890",
    "deleted_at": null
}).pretty()

# Создание индекса для поиска по email
db.users.createIndex({"email": 1}, {unique: true})

# Просмотр всех индексов
db.users.getIndexes()

# Статистика коллекции
db.posts.stats()

# Удаление документа по _id
db.posts.deleteOne({"_id": ObjectId("67f4b8a9d2f4a12c34567891")})

# Агрегация: количество постов по статусам
db.posts.aggregate([
    {"$match": {"deleted_at": null}},
    {"$group": {"_id": "$status", "count": {"$sum": 1}}}
])</code></pre>

  <h3>Тестирование работы API с MongoDB</h3>
  
  <h4>1. Полный цикл CRUD операций</h4>
  <pre><code># Регистрация
curl -X POST http://localhost:4200/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@mongodb.com","password":"test123"}' | jq .

# Логин и сохранение cookies
curl -X POST http://localhost:4200/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@mongodb.com","password":"test123"}' \
  -c cookies.txt | jq .

# Получение ID пользователя
USER_ID=$(curl -s http://localhost:4200/api/v1/auth/whoami -b cookies.txt | jq -r '.data.id')

# Создание поста
curl -X POST http://localhost:4200/api/v1/posts \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"title\":\"MongoDB Post\",\"status\":\"active\"}" \
  -b cookies.txt | jq .

# Получение всех постов (с кешированием)
time curl -s "http://localhost:4200/api/v1/posts?page=1&limit=10" -b cookies.txt > /dev/null

# Обновление поста
POST_ID=$(curl -s "http://localhost:4200/api/v1/posts?page=1&limit=10" -b cookies.txt | jq -r '.data.data[0].id')
curl -X PUT "http://localhost:4200/api/v1/posts/$POST_ID" \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title"}' \
  -b cookies.txt | jq .

# Удаление поста
curl -X DELETE "http://localhost:4200/api/v1/posts/$POST_ID" -b cookies.txt

# Выход (инвалидация JTI)
curl -X POST http://localhost:4200/api/v1/auth/logout -b cookies.txt</code></pre>

  <h4>2. Проверка в MongoDB Express</h4>
  <pre><code># Откройте браузер и перейдите по адресу:
http://localhost:8081
  
  <p>Скрипт для переноса данных из PostgreSQL в MongoDB:</p>
  <pre><code>// scripts/migrate.go
package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "log"
    
    _ "github.com/lib/pq"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
    ID       string `json:"id"`
    Email    string `json:"email"`
    Phone    string `json:"phone"`
    CreatedAt string `json:"created_at"`
}

func main() {
    // Подключение к PostgreSQL
    pgDB, _ := sql.Open("postgres", "postgres://user:pass@localhost:5432/rugram_db")
    
    // Подключение к MongoDB
    mongoClient, _ := mongo.Connect(context.Background(), 
        options.Client().ApplyURI("mongodb://localhost:27017"))
    mongoDB := mongoClient.Database("rugram_db")
    
    // Чтение из PostgreSQL
    rows, _ := pgDB.Query("SELECT id, email, phone, created_at FROM users WHERE deleted_at IS NULL")
    
    // Вставка в MongoDB
    for rows.Next() {
        var user User
        rows.Scan(&user.ID, &user.Email, &user.Phone, &user.CreatedAt)
        
        // Конвертация UUID в ObjectID
        data, _ := json.Marshal(user)
        var mongoDoc map[string]interface{}
        json.Unmarshal(data, &mongoDoc)
        
        mongoDB.Collection("users").InsertOne(context.Background(), mongoDoc)
    }
    
    log.Println("Migration completed!")
}</code></pre>

  <h3>Полезные команды Docker</h3>
  
  <pre><code># Запуск всех сервисов
docker-compose up -d

# Просмотр логов MongoDB
docker logs rugram_mongo -f

# Просмотр логов приложения
docker logs rugram_api -f

# Подключение к MongoDB Shell
docker exec -it rugram_mongo mongosh -u rugram_user -p rugram_password

# Бэкап MongoDB
docker exec rugram_mongo mongodump --username rugram_user --password rugram_password \
  --authenticationDatabase admin --db rugram_db --out /dump

# Восстановление из бэкапа
docker exec rugram_mongo mongorestore --username rugram_user --password rugram_password \
  --authenticationDatabase admin --db rugram_db /dump/rugram_db

# Очистка всех данных (осторожно!)
docker-compose down -v

# Пересборка без кеша
docker-compose build --no-cache</code></pre>

  <h3>Устранение неполадок</h3>
  
  <h4>Проблема: Не создаются коллекции в MongoDB</h4>
  <pre><code># Решение: коллекции создаются автоматически при первой вставке данных
# Выполните регистрацию пользователя:
curl -X POST http://localhost:4200/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}'

# После этого в MongoDB Express появятся коллекции users и user_tokens</code></pre>

  <h4>Проблема: Ошибка подключения к MongoDB</h4>
  <pre><code># Проверьте статус контейнера
docker ps | grep mongo

# Проверьте логи
docker logs rugram_mongo

# Перезапустите сервисы
docker-compose restart mongo app</code></pre>

  <h4>Проблема: Медленные запросы</h4>
  <pre><code># Проверьте созданы ли индексы
docker exec -it rugram_mongo mongosh -u rugram_user -p rugram_password
use rugram_db
db.users.getIndexes()
db.posts.getIndexes()

# Создайте недостающие индексы
db.users.createIndex({"email": 1}, {unique: true})
db.posts.createIndex({"user_id": 1, "created_at": -1})</code></pre>

  <h3>Выводы по лабораторной работе</h3>
  
  <ul>
    <li><b>MongoDB обеспечивает более высокую производительность</b> для операций чтения (на 40-50%) благодаря отсутствию JOIN и нативному хранению JSON-подобных документов</li>
    <li><b>Гибкая схема</b> позволяет легко добавлять новые поля без миграций</li>
    <li><b>Встроенные массивы</b> идеально подходят для хранения комментариев, лайков и других связанных данных</li>
    <li><b>Простота масштабирования</b> - шардинг в MongoDB проще чем партиционирование в PostgreSQL</li>
    <li><b>Mongo Express</b> предоставляет удобный веб-интерфейс для администрирования</li>
    <li><b>Сохранилась полная совместимость</b> с Redis кешированием и JWT аутентификацией</li>
  </ul>
  
  <p><b>Итог:</b> Для социальной сети с большим количеством операций чтения и не очень сложными связями MongoDB является оптимальным выбором. PostgreSQL остается лучшим выбором для систем с критичными транзакциями и сложными аналитическими запросами.</p>

</details>