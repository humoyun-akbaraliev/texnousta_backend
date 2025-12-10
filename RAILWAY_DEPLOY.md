# 🚀 Деплой на Railway (рекомендуется для Go)

## Быстрый деплой (3 команды):
```bash
npm install -g @railway/cli
railway login
railway init
railway up
```

## Файлы уже готовы:
✅ railway.toml - есть
✅ Dockerfile - есть  
✅ .env - настроен
✅ Procfile - есть

## После деплоя:
- API: https://your-app.railway.app
- Swagger: https://your-app.railway.app/swagger/index.html

## Переменные окружения в Railway:
GIN_MODE=release
DATABASE_URL=postgresql://... (автоматически из Railway PostgreSQL)

## Преимущества Railway vs Vercel:
✅ Поддержка Go
✅ PostgreSQL база данных
✅ Автодеплой из Git  
✅ Бесплатный план
✅ Логи и мониторинг