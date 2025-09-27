# 🚀 AIEN Backend - Руководство по развертыванию

## 📋 Обзор

Это руководство поможет вам развернуть AIEN Backend на любой машине с помощью Docker Compose.

## 🖥️ Системные требования

### Минимальные требования:
- **OS:** Linux (Ubuntu 20.04+), macOS (10.15+), Windows 10+
- **RAM:** 8GB (рекомендуется 16GB)
- **CPU:** 4 ядра (рекомендуется 8 ядер)
- **Диск:** 20GB свободного места
- **Docker:** 20.10+
- **Docker Compose:** 2.0+

### Рекомендуемые требования:
- **RAM:** 16GB+ (для работы с большими файлами)
- **GPU:** NVIDIA GPU с поддержкой CUDA (для Ollama)
- **Диск:** SSD 50GB+ (для быстрой работы)

## 🔧 Установка зависимостей

### Ubuntu/Debian:
```bash
# Обновление системы
sudo apt update && sudo apt upgrade -y

# Установка Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Установка Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Перезагрузка для применения изменений
sudo reboot
```

### CentOS/RHEL:
```bash
# Установка Docker
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install -y docker-ce docker-ce-cli containerd.io
sudo systemctl start docker
sudo systemctl enable docker

# Установка Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### macOS:
```bash
# Установка Homebrew (если не установлен)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Установка Docker Desktop
brew install --cask docker
```

### Windows:
1. Скачайте Docker Desktop: https://www.docker.com/products/docker-desktop
2. Установите и запустите Docker Desktop
3. Включите WSL 2 (рекомендуется)

## 📥 Развертывание

### 1. Клонирование репозитория:
```bash
git clone https://github.com/your-username/AIEN_backend.git
cd AIEN_backend
```

### 2. Проверка зависимостей:
```bash
# Проверить Docker
docker --version
docker-compose --version

# Проверить доступность портов
./start.sh status
```

### 3. Запуск приложения:
```bash
# Запуск всех сервисов
./start.sh start
```

### 4. Проверка работоспособности:
```bash
# Проверить статус
./start.sh status

# Проверить логи
./start.sh logs
```

## 🔧 Конфигурация

### Настройка портов:
Если порты заняты, измените их в `docker-compose.prod.yml`:

```yaml
services:
  api-gateway:
    ports:
      - "8080:8080"  # Измените на другой порт
```

### Настройка ресурсов:
Для ограничения ресурсов добавьте в `docker-compose.prod.yml`:

```yaml
services:
  ollama:
    deploy:
      resources:
        limits:
          memory: 4G
          cpus: '2.0'
```

### Настройка GPU:
Для использования GPU с Ollama:

```yaml
services:
  ollama:
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]
```

## 🧪 Тестирование

### 1. Проверка здоровья сервисов:
```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8083/health
```

### 2. Тестирование загрузки файла:
```bash
curl -X POST -F "file=@test_data/sales_data.csv" "http://localhost:8081/api/v1/files/upload?user_id=test-user"
```

### 3. Тестирование анализа:
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"file_id":"test-file","user_id":"test-user","file_path":"path/to/file"}' \
  http://localhost:8083/api/v1/analysis/start
```

## 📊 Мониторинг

### Просмотр логов:
```bash
# Все сервисы
./start.sh logs

# Конкретный сервис
docker-compose -f docker-compose.prod.yml logs api-gateway
```

### Мониторинг ресурсов:
```bash
# Использование ресурсов
docker stats

# Использование диска
docker system df
```

### Проверка статуса:
```bash
# Статус контейнеров
docker-compose -f docker-compose.prod.yml ps

# Статус сервисов
./start.sh status
```

## 🗄️ Управление данными

### Резервное копирование:
```bash
# Создать бэкап базы данных
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U postgres aien_db > backup.sql

# Создать бэкап MinIO
docker-compose -f docker-compose.prod.yml exec minio mc mirror /data /backup
```

### Восстановление:
```bash
# Восстановить базу данных
docker-compose -f docker-compose.prod.yml exec -T postgres psql -U postgres aien_db < backup.sql
```

### Очистка данных:
```bash
# Очистить все данные
./start.sh clean
```

## 🚨 Устранение неполадок

### Проблемы с портами:
```bash
# Найти процесс, использующий порт
sudo lsof -i :8080

# Убить процесс
sudo kill -9 <PID>
```

### Проблемы с памятью:
```bash
# Проверить использование памяти
free -h
docker stats

# Очистить Docker
docker system prune -a
```

### Проблемы с GPU:
```bash
# Проверить NVIDIA Docker
docker run --rm --gpus all nvidia/cuda:11.0-base nvidia-smi

# Установить NVIDIA Docker (Ubuntu)
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list
sudo apt-get update && sudo apt-get install -y nvidia-docker2
sudo systemctl restart docker
```

### Проблемы с сетью:
```bash
# Проверить сеть Docker
docker network ls
docker network inspect aien_backend_aien_network

# Пересоздать сеть
docker-compose -f docker-compose.prod.yml down
docker network prune
./start.sh start
```

## 🔒 Безопасность

### Настройка файрвола:
```bash
# Ubuntu/Debian
sudo ufw allow 8080
sudo ufw allow 8081
sudo ufw allow 8083
sudo ufw enable
```

### Изменение паролей:
```bash
# Изменить пароли в docker-compose.prod.yml
POSTGRES_PASSWORD: your_secure_password
MINIO_ROOT_PASSWORD: your_secure_password
```

### SSL/TLS:
Для продакшена настройте reverse proxy (nginx) с SSL сертификатами.

## 📈 Масштабирование

### Горизонтальное масштабирование:
```bash
# Увеличить количество реплик
docker-compose -f docker-compose.prod.yml up -d --scale api-gateway=3
```

### Вертикальное масштабирование:
Измените ресурсы в `docker-compose.prod.yml`:

```yaml
services:
  api-gateway:
    deploy:
      resources:
        limits:
          memory: 2G
          cpus: '1.0'
```

## 🎯 Продакшен

### Переменные окружения:
Создайте `.env` файл:

```bash
POSTGRES_PASSWORD=your_secure_password
MINIO_ROOT_PASSWORD=your_secure_password
REDIS_PASSWORD=your_secure_password
```

### Мониторинг:
- Настройте Prometheus + Grafana
- Настройте ELK Stack для логов
- Настройте AlertManager для уведомлений

### Резервное копирование:
- Настройте автоматические бэкапы
- Настройте репликацию базы данных
- Настройте синхронизацию с облачным хранилищем

## 📞 Поддержка

При возникновении проблем:

1. Проверьте логи: `./start.sh logs`
2. Проверьте статус: `./start.sh status`
3. Создайте issue в GitHub
4. Обратитесь к документации

---

**AIEN Backend** - Готов к развертыванию! 🚀