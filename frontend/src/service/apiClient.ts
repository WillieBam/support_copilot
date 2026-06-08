import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

let isRefreshing = false;
let failedQueue: Array<{
  resolve: () => void; 
  reject: (error: AxiosError) => void;
}> = [];

const processQueue = (error: AxiosError | null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(); 
    }
  });
  failedQueue = [];
};

// handle 401 and silent refresh
apiClient.interceptors.response.use(
  (response) => response, 
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      
      if (isRefreshing) {
        return new Promise<void>(function(resolve, reject) {
          failedQueue.push({ resolve, reject });
        })
        .then(() => {return apiClient(originalRequest);}) // the browser already has the new cookie, just resend the request
        .catch((err) => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // call backend to get a new access cookie
        // browser will automatically send the HttpOnly refresh cookie with this request
        await axios.post(`${apiClient.defaults.baseURL}/auth/refresh`, {}, {
          withCredentials: true // explicitly allow cookies
        });

        // if successful, backend set a NEW HttpOnly access cookie.
        // process the queue, resolving all paused requests
        processQueue(null);
        return apiClient(originalRequest);

      } catch (refreshError) {
        // if refresh fails, force logout
        processQueue(refreshError as AxiosError);
        await axios.post(`${apiClient.defaults.baseURL}/auth/logout`, {}, { withCredentials: true }).catch(() => {});
        
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