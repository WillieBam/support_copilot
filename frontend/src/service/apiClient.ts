import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios';
import { firebaseAuth } from '@/firebase';
import { exchangeToken } from './auth/authService';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

let isRefreshing = false;
let failedQueue: Array<{
  resolve:(token: string) => void; 
  reject: (error: AxiosError) => void;
}> = [];

const processQueue = (error: AxiosError | null, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else if (token) {
      prom.resolve(token); 
    }
  });
  failedQueue = [];
};
// handle 401 and silent refresh
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('support_copilot_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

apiClient.interceptors.response.use(
  (response) => response, 
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      // If a token refresh cycle is already ongoing, queue this request
      if (isRefreshing) {
        return new Promise<string>(function(resolve, reject) {
          failedQueue.push({ resolve, reject });
        })
      .then((newToken) => {
          // Update authorization header with the fresh token issued while this request waited
          originalRequest.headers.Authorization = `Bearer ${newToken}`;
          return apiClient(originalRequest);
        })
      .catch((err) => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const user = firebaseAuth.currentUser;
        if (!user) {
          throw new Error("No active Firebase session found during silent refresh");
        }
        const freshBackendToken = await exchangeToken(user);
        processQueue(null, freshBackendToken);
        originalRequest.headers.Authorization = `Bearer ${freshBackendToken}`;
        return apiClient(originalRequest);

      } catch (refreshError) {
        // if refresh fails, force logout
        processQueue(refreshError as AxiosError, null);
        localStorage.removeItem('support_copilot_token');
        await firebaseAuth.signOut().catch(() => {});
        window.location.href = '/login';
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;