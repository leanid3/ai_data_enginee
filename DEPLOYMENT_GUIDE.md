# 🚀 AIEN Backend - Инструкция по развертыванию

## 📦 Установка

### 1. Распаковка архива
```bash
tar -xzf AIEN_Backend_*.tar.gz
cd AIEN_backend
```

### 2. Установка зависимостей
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y docker.io docker-compose curl

# CentOS/RHEL
sudo yum install -y docker docker-compose curl

# macOS
brew install docker docker-compose curl
```

### 3. Запуск системы
```bash
chmod +x *.sh
./start.sh
```

## 🔧 Настройка

### Порты (по умолчанию)
- 8080 - API Gateway
- 8081 - File Service  
- 8083 - Data Analysis Service
- 8084 - Adminer
- 9000/9001 - MinIO
- 11434 - Ollama

### Изменение портов
Отредактируйте `docker-compose.prod.yml`:
```yaml
ports:
  - "8080:8080"  # Измените на нужный порт
```

## 🧪 Тестирование

### 1. Проверка здоровья
```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8083/health
```

### 2. Postman коллекция
1. Импортируйте `AIEN_Backend_API.postman_collection.json`
2. Настройте переменные
3. Запустите тесты

### 3. Тестовые файлы
Используйте файлы из папки `test_data/`:
- `sales_data.csv`
- `users.json`
- `products.xml`

## 📊 Мониторинг

### Логи
```bash
docker-compose logs -f
docker-compose logs -f [service]
```

### Статистика
```bash
docker stats
docker system df
```

## 🔧 Управление

### Запуск
```bash
./start.sh
```

### Остановка
```bash
./stop.sh
```

### Перезапуск
```bash
./restart.sh
```

### Очистка
```bash
./clean.sh
```

## 🐛 Решение проблем

### Порты заняты
```bash
lsof -i :8080
sudo kill -9 [PID]
```

### Недостаточно памяти
```bash
docker system prune -f
```

### Медленная работа
```bash
docker stats
```

## 📞 Поддержка

### Полезные команды
```bash
# Очистка системы
./clean.sh

# Перезапуск
./restart.sh

# Просмотр логов
docker-compose logs -f

# Проверка статуса
docker-compose ps
```

### Частые проблемы
1. **Порты заняты** - проверьте `lsof -i :PORT`
2. **Недостаточно памяти** - увеличьте RAM
3. **Медленная работа** - проверьте ресурсы

## 📄 Лицензия

MIT License

## 👥 Команда

AIEN Team

---

**Готово к использованию!** 🎉
