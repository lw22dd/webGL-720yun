import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";

const serverConfig = {
    baseURL: "/api/v1",
    useTokenAuthorization: true,
};

let isRefreshing = false;
let refreshSubscribers: ((token: string) => void)[] = [];

const subscribeTokenRefresh = (callback: (token: string) => void) => {
    refreshSubscribers.push(callback);
};

const onTokenRefreshed = (token: string) => {
    refreshSubscribers.forEach(callback => callback(token));
    refreshSubscribers = [];
};

const serviceAxios = axios.create({
    baseURL: serverConfig.baseURL,
    timeout: 10000,
    withCredentials: false,
});

const getToken = () => {
    try {
        const userData = localStorage.getItem('user');
        if (userData) {
            const parsed = JSON.parse(userData);
            return parsed.accessToken || '';
        }
    } catch (e) {
        console.error('Failed to get token from localStorage:', e);
    }
    return '';
};

const getRefreshToken = () => {
    try {
        const userData = localStorage.getItem('user');
        if (userData) {
            const parsed = JSON.parse(userData);
            return parsed.refreshToken || '';
        }
    } catch (e) {
        console.error('Failed to get refresh token from localStorage:', e);
    }
    return '';
};

const setToken = (accessToken: string, refreshToken: string) => {
    try {
        const userData = localStorage.getItem('user');
        if (userData) {
            const parsed = JSON.parse(userData);
            parsed.accessToken = accessToken;
            parsed.refreshToken = refreshToken;
            localStorage.setItem('user', JSON.stringify(parsed));
        }
    } catch (e) {
        console.error('Failed to set token in localStorage:', e);
    }
};

const removeToken = () => {
    try {
        const userData = localStorage.getItem('user');
        if (userData) {
            const parsed = JSON.parse(userData);
            parsed.accessToken = '';
            parsed.refreshToken = '';
            localStorage.setItem('user', JSON.stringify(parsed));
        }
    } catch (e) {
        console.error('Failed to remove token from localStorage:', e);
    }
};

const refreshToken = async (): Promise<string | null> => {
    const refreshTokenValue = getRefreshToken();
    if (!refreshTokenValue) {
        return null;
    }

    try {
        const response = await axios.post('/api/v1/auth/refresh', {
            refresh_token: refreshTokenValue
        });

        if (response.data.code === 200 && response.data.data) {
            const { access_token } = response.data.data;
            setToken(access_token, refreshTokenValue);
            return access_token;
        }
        return null;
    } catch (error) {
        console.error('Token refresh failed:', error);
        return null;
    }
};

const resetStoreAndRedirect = () => {
    removeToken();
    localStorage.removeItem('user');
    window.location.href = '/';
};

serviceAxios.interceptors.request.use(
    (config) => {
        if (config.method?.toLowerCase() !== 'get' && !config.headers['Content-Type']) {
            if (config.data && !(config.data instanceof FormData)) {
                config.headers['Content-Type'] = 'application/json';
            }
        }

        const token = getToken();
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

serviceAxios.interceptors.response.use(
    (res) => {
        let data = res.data;
        if (data.code && data.code !== 200) {
            return Promise.reject(new Error(data.msg || '请求失败'));
        }
        return data;
    },
    async (err: AxiosError) => {
        const originalRequest = err.config as InternalAxiosRequestConfig & { _retry?: boolean };
        const isLoginRequest = originalRequest.url?.includes('/auth/login');

        if (err.response?.status === 401 && !originalRequest._retry && !isLoginRequest) {
            if (isRefreshing) {
                return new Promise((resolve) => {
                    subscribeTokenRefresh((token: string) => {
                        if (originalRequest.headers) {
                            originalRequest.headers['Authorization'] = `Bearer ${token}`;
                        }
                        resolve(serviceAxios(originalRequest));
                    });
                });
            }

            originalRequest._retry = true;
            isRefreshing = true;

            try {
                const newToken = await refreshToken();
                if (newToken) {
                    isRefreshing = false;
                    onTokenRefreshed(newToken);
                    if (originalRequest.headers) {
                        originalRequest.headers['Authorization'] = `Bearer ${newToken}`;
                    }
                    return serviceAxios(originalRequest);
                } else {
                    isRefreshing = false;
                    resetStoreAndRedirect();
                    return Promise.reject(new Error('Token refresh failed'));
                }
            } catch (refreshError) {
                isRefreshing = false;
                resetStoreAndRedirect();
                return Promise.reject(refreshError);
            }
        }

        if (err.response) {
            const { status, data } = err.response as { status: number; data: any };
            let message = data?.msg || '请求失败';
            if (status === 401) {
                message = data?.msg || '登录已过期，请重新登录';
            } else if (status === 400) {
                message = data?.msg || '请求参数错误';
            } else if (status === 404) {
                message = data?.msg || '资源不存在';
            } else if (status === 500) {
                message = data?.msg || '服务器内部错误';
            }
            return Promise.reject(new Error(message));
        }
        return Promise.reject(err);
    }
);

export default serviceAxios;
