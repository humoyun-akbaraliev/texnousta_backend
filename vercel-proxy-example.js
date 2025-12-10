# Vercel НЕ ПОДДЕРЖИВАЕТ Go серверы!
# Используйте Railway или Render вместо этого.

# Если все же нужен Vercel, то только через proxy:
# 1. Деплой Go API на Railway/Render
# 2. Создать Next.js приложение на Vercel как прокси

# Пример Next.js API route для проксирования:
# pages/api/[...slug].js

export default async function handler(req, res) {
  const API_URL = 'https://your-railway-app.railway.app'
  const url = `${API_URL}/api/v1/${req.query.slug.join('/')}`
  
  try {
    const response = await fetch(url, {
      method: req.method,
      headers: {
        'Content-Type': 'application/json',
        ...req.headers
      },
      body: req.method !== 'GET' ? JSON.stringify(req.body) : undefined
    })
    
    const data = await response.json()
    res.status(response.status).json(data)
  } catch (error) {
    res.status(500).json({ error: 'Proxy error' })
  }
}