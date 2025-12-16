# 第一阶段：构建前端应用
FROM node:18-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json
COPY package.json package-lock.json ./

# 安装依赖
RUN npm ci

# 复制所有代码
COPY . .

# 构建应用
RUN npm run build

# 第二阶段：运行应用
FROM nginx:alpine

# 复制构建好的静态文件到nginx的html目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 设置静态文件权限
RUN chown -R nginx:nginx /usr/share/nginx/html && \ 
    chmod -R 755 /usr/share/nginx/html

# 复制nginx配置文件
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 暴露端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]