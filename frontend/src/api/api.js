import axios from 'axios';

// Create an Axios instance with the base URL of your backend
const api = axios.create({
  baseURL: 'http://localhost:8080', // Adjust this to match your backend server URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add the Authorization header
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token'); // Use 'access_token' consistently
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor to handle token refresh on 401 errors
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      try {
        const refreshToken = localStorage.getItem('refresh_token');
        const response = await api.post('/refresh', { refresh_token: refreshToken });
        const { access_token } = response.data;
        localStorage.setItem('access_token', access_token); // Store as 'access_token'
        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return api(originalRequest); // Retry the original request
      } catch (refreshError) {
        // If refresh fails, clear tokens and redirect to login
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        window.location.href = '/login';
        return Promise.reject(refreshError);
      }
    }
    return Promise.reject(error);
  }
);

export default api;