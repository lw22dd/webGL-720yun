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
    // 处理自己的业务逻辑，比如判断 token 是否过期等等
    // 代码块
    return data;
  },
  (err) => {
    // 处理请求错误，这里可以用全局提示框进行提示
    // 错误处理
    return Promise.reject(err);
  }
);

export default serviceAxios;
