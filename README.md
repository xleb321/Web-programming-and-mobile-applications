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

<details open>
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
      <td>❌</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/posts/:id</td>
      <td>Пост по ID</td>
      <td>❌</td>
    </tr>
    <tr>
      <td>POST</td>
      <td>/api/v1/posts</td>
      <td>Создать пост (привязан к user_id)</td>
      <td>✅</td>
    </tr>
    <tr>
      <td>PUT/PATCH</td>
      <td>/api/v1/posts/:id</td>
      <td>Обновить свой пост</td>
      <td>✅</td>
    </tr>
    <tr>
      <td>DELETE</td>
      <td>/api/v1/posts/:id</td>
      <td>Удалить свой пост</td>
      <td>✅</td>
    </tr>
    <tr>
      <td>GET</td>
      <td>/api/v1/posts/user/:userId</td>
      <td>Посты пользователя</td>
      <td>❌</td>
    </tr>
  </table>
  
  <h3>📝 Примеры запросов</h3>
  
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