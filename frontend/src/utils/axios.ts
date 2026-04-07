import axios from "axios";
import { useUserStore } from "@/stores/userStore";

const serverConfig = {
    baseURL: "/api/v1", // 请求基础地址,使用相对路径，通过nginx反向代理到后端
    useTokenAuthorization: false, // 是否开启 token 认证
  };
// 创建 axios 请求实例
const serviceAxios = axios.create({
  baseURL: serverConfig.baseURL, // 基础请求地址
  timeout: 10000, // 请求超时设置
  withCredentials: false, // 跨域请求是否需要携带 cookie
});

// 请求拦截器
serviceAxios.interceptors.request.use(
  (config) => {
    // 确保请求头为application/json
    if (config.method?.toLowerCase() !== 'get' && !config.headers['Content-Type']) {
      config.headers['Content-Type'] = 'application/json';
    }
    
    // 添加认证token，从store获取
    const userStore = useUserStore();
    const token = userStore.accessToken;
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
serviceAxios.interceptors.response.use(
  (res) => {
    let data = res.data;
    if (data.code && data.code !== 200) {
      return Promise.reject(new Error(data.msg || '请求失败'));
    }
    return data;
  },
  (err) => {
    if (err.response) {
      const { status, data } = err.response;
      let message = data?.msg || '请求失败';
      if (status === 401) {
        message = data?.msg || '用户名或密码错误';
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
