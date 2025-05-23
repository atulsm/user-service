import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Add response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
}

export interface User {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateUserRequest {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  phoneNumber?: string;
}

export interface UpdateUserRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
  phoneNumber?: string;
}

interface UserActivity {
  date: string;
  newUsers: number;
  activeUsers: number;
}

interface UserStats {
  totalUsers: number;
  activeUsers: number;
  newUsers: number;
}

interface UserActivityParams {
  startDate: string;
  endDate: string;
}

export const authService = {
  login: async (data: LoginRequest) => {
    const response = await api.post('/auth/login', data);
    const { token } = response.data;
    localStorage.setItem('token', token);
    return response.data;
  },

  register: async (data: RegisterRequest) => {
    const response = await api.post('/auth/register', data);
    return response.data;
  },

  logout: async () => {
    try {
      await api.post('/auth/logout');
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      localStorage.removeItem('token');
    }
  },
};

export const userService = {
  getUsers: () => api.get<User[]>('/users').then((res) => res.data),
  getUser: (id: string) => api.get<User>(`/users/${id}`).then((res) => res.data),
  createUser: (data: CreateUserRequest) => api.post<User>('/users', data).then((res) => res.data),
  updateUser: (id: string, data: UpdateUserRequest) =>
    api.put<User>(`/users/${id}`, data).then((res) => res.data),
  deleteUser: (id: string) => api.delete(`/users/${id}`).then((res) => res.data),
  getProfile: () => api.get<User>('/users/profile').then((res) => res.data),
  updateProfile: (data: UpdateUserRequest) =>
    api.put<User>('/users/profile', data).then((res) => res.data),
  getUserActivity: (params: UserActivityParams) => api.get<UserActivity[]>('/users/activity', { params }),
  getUserStats: () => api.get<UserStats>('/users/stats'),
}; 