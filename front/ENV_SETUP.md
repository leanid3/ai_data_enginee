# Настройка переменных окружения

## Создание .env файла

Создайте файл `.env` в корне проекта frontend со следующим содержимым:

```env
# API Configuration
REACT_APP_API_BASE_URL=http://localhost:8080
REACT_APP_USER_ID=test-user-test

# Development settings
REACT_APP_DEBUG=true
```

## Описание переменных

- `REACT_APP_API_BASE_URL` - базовый URL для API backend
- `REACT_APP_USER_ID` - идентификатор пользователя для API запросов
- `REACT_APP_DEBUG` - включение режима отладки

## Настройка для разных окружений

### Development
```env
REACT_APP_API_BASE_URL=http://localhost:8080
REACT_APP_USER_ID=dev-user
REACT_APP_DEBUG=true
```

### Production
```env
REACT_APP_API_BASE_URL=https://api.yourdomain.com
REACT_APP_USER_ID=prod-user
REACT_APP_DEBUG=false
```

### Staging
```env
REACT_APP_API_BASE_URL=https://staging-api.yourdomain.com
REACT_APP_USER_ID=staging-user
REACT_APP_DEBUG=true
```

## Важные замечания

1. Все переменные окружения в React должны начинаться с `REACT_APP_`
2. После изменения .env файла необходимо перезапустить сервер разработки
3. Не добавляйте .env файл в систему контроля версий (он уже в .gitignore)
4. Для продакшена настройте переменные окружения на сервере
