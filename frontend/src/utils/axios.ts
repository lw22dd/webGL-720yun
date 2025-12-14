
import axios from "axios";

const serverConfig = {
    baseURL: "/api", // 请求基础地址,使用相对路径，通过nginx反向代理到后端
    useTokenAuthorization: false, // 是否开启 token 认证
  };
// 创建 axios 请求实例
const serviceAxios = axios.create({
  baseURL: serverConfig.baseURL, // 基础请求地址
  timeout: 10000, // 请求超时设置
  withCredentials: false, // 跨域请求是否需要携带 cookie
});
//拦截器可选
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