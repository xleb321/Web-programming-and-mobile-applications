Write-Host "=== RuGramm Instagram Clone ===" -ForegroundColor Cyan

# Остановка старых контейнеров
Write-Host "`nОстанавливаем старые контейнеры..." -ForegroundColor Yellow
docker-compose down -v

# Очистка Docker
Write-Host "`nОчищаем Docker кэш..." -ForegroundColor Yellow
docker system prune -f

# Сборка и запуск
Write-Host "`nСобираем и запускаем контейнеры..." -ForegroundColor Green
docker-compose up --build