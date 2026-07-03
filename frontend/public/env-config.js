// Fallback local — sobrescrito em produção pelo start.sh
window.__ENV__ = window.__ENV__ || {
  API_URL: '/api/v1',
}
