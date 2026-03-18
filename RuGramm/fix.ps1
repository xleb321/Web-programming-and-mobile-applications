Write-Host "=== Полное исправление Docker сборки ===" -ForegroundColor Cyan

# Шаг 1: Остановка всех контейнеров и очистка
Write-Host "`n1. Останавливаем контейнеры и очищаем Docker..." -ForegroundColor Yellow
docker-compose down -v
docker system prune -af

# Шаг 2: Удаление старых модулей
Write-Host "`n2. Удаляем старые модули Go..." -ForegroundColor Yellow
Remove-Item go.mod, go.sum -Force -ErrorAction SilentlyContinue

# Шаг 3: Инициализация нового модуля
Write-Host "`n3. Инициализируем новый модуль..." -ForegroundColor Green
go mod init RuGramm

# Шаг 4: Установка зависимостей
Write-Host "`n4. Устанавливаем зависимости..." -ForegroundColor Green
$dependencies = @(
    "github.com/gin-gonic/gin@v1.9.1",
    "github.com/go-playground/validator/v10@v10.16.0",
    "github.com/google/uuid@v1.4.0",
    "github.com/joho/godotenv@v1.5.1",
    "gorm.io/gorm@v1.25.5",
    "gorm.io/driver/postgres@v1.5.4"
)

foreach ($dep in $dependencies) {
    Write-Host "   Installing $dep..." -ForegroundColor Gray
    go get $dep
}

# Шаг 5: Синхронизация
Write-Host "" -ForegroundColor Green
go mod tidy

# Шаг 6: Проверка go.mod
Write-Host "" -ForegroundColor Cyan
Get-Content go.mod

# Шаг 7: Сборка Docker
Write-Host "" -ForegroundColor Green
docker-compose build --no-cache

# Шаг 8: Запуск
Write-Host "" -ForegroundColor Green
docker-compose up

Write-Host "" -ForegroundColor Green