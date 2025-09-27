# 🚀 AIEN Backend - Быстрая установка

## ⚡ Быстрый старт (3 команды)

```bash
# 1. Клонировать репозиторий
git clone https://github.com/your-username/AIEN_backend.git
cd AIEN_backend

# 2. Запустить приложение
./start.sh

# 3. Проверить статус
./start.sh status
```

## 🎯 Готово!

Ваше приложение запущено и доступно по адресам:

- **🌐 API Gateway:** http://localhost:8080
- **📁 File Service:** http://localhost:8081  
- **🔍 Data Analysis:** http://localhost:8083
- **🗄️ Database Admin:** http://localhost:8084
- **📦 MinIO Console:** http://localhost:9001

## 🧪 Тестирование

### Загрузить файл:
```bash
curl -X POST -F "file=@test_data/sales_data.csv" "http://localhost:8081/api/v1/files/upload?user_id=test-user"
```

### Запустить анализ:
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"file_id":"test-file","user_id":"test-user","file_path":"path/to/file"}' \
  http://localhost:8083/api/v1/analysis/start
```

## 📊 Postman коллекция

1. Импортируйте `AIEN_Backend_API.postman_collection.json` в Postman
2. Настройте пути к файлам в переменных
3. Запустите тестовые сценарии

## 🛠️ Управление

```bash
./start.sh start     # Запустить
./start.sh stop      # Остановить  
./start.sh restart   # Перезапустить
./start.sh status    # Статус
./start.sh logs      # Логи
./start.sh clean     # Очистить данные
```

## 📋 Требования

- Docker 20.10+
- Docker Compose 2.0+
- 8GB RAM
- 20GB свободного места

## 🆘 Помощь

- **Документация:** README.md
- **Развертывание:** DEPLOYMENT.md
- **Тестирование:** POSTMAN_TESTING_GUIDE.md
- **Issues:** GitHub Issues

---

**AIEN Backend** - Система анализа данных с LLM 🚀