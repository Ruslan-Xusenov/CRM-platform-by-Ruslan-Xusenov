import axios from 'axios';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: `${API_BASE}/api/v1`,
  headers: { 'Content-Type': 'application/json' },
});

// Inject access token into every request
api.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('access_token');
    if (token) config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auto-refresh on 401
api.interceptors.response.use(
  (res) => res,
  async (error) => {
    const original = error.config;
    if (error.response?.status === 401 && !original._retry) {
      original._retry = true;
      try {
        const refreshToken = localStorage.getItem('refresh_token');
        const res = await axios.post(`${API_BASE}/api/v1/auth/refresh`, { refresh_token: refreshToken });
        localStorage.setItem('access_token', res.data.access_token);
        localStorage.setItem('refresh_token', res.data.refresh_token);
        original.headers.Authorization = `Bearer ${res.data.access_token}`;
        return api(original);
      } catch {
        localStorage.clear();
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

export default api;

// ─── Auth API ────────────────────────────────────────────────
export const authAPI = {
  login: (data: { email: string; password: string }) => api.post('/auth/login', data),
  register: (data: { email: string; password: string; first_name: string; last_name: string; tenant_name: string }) => api.post('/auth/register', data),
  logout: (refresh_token: string) => api.post('/auth/logout', { refresh_token }),
};

// ─── CRM API ─────────────────────────────────────────────────
export const leadsAPI = {
  list: (params?: Record<string, string>) => api.get('/leads', { params }),
  get: (id: string) => api.get(`/leads/${id}`),
  create: (data: Record<string, unknown>) => api.post('/leads', data),
  update: (id: string, data: Record<string, unknown>) => api.put(`/leads/${id}`, data),
  delete: (id: string) => api.delete(`/leads/${id}`),
  convert: (id: string) => api.post(`/leads/${id}/convert`),
};

export const contactsAPI = {
  list: (params?: Record<string, string>) => api.get('/contacts', { params }),
  get: (id: string) => api.get(`/contacts/${id}`),
  create: (data: Record<string, unknown>) => api.post('/contacts', data),
  update: (id: string, data: Record<string, unknown>) => api.put(`/contacts/${id}`, data),
  delete: (id: string) => api.delete(`/contacts/${id}`),
};

export const dealsAPI = {
  list: (params?: Record<string, string>) => api.get('/deals', { params }),
  get: (id: string) => api.get(`/deals/${id}`),
  create: (data: Record<string, unknown>) => api.post('/deals', data),
  update: (id: string, data: Record<string, unknown>) => api.put(`/deals/${id}`, data),
  delete: (id: string) => api.delete(`/deals/${id}`),
};

export const pipelinesAPI = {
  list: () => api.get('/pipelines'),
  create: (data: Record<string, unknown>) => api.post('/pipelines', data),
};

// ─── PBX API ─────────────────────────────────────────────────
export const callsAPI = {
  active: () => api.get('/calls/active'),
  history: (params?: Record<string, string>) => api.get('/calls/history', { params }),
  originate: (to: string) => api.post('/calls/originate', { to }),
};

export const extensionsAPI = {
  list: () => api.get('/pbx/extensions'),
};
