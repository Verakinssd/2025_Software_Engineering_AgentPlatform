
# AI 智能低代码开发平台

## 运行说明

本项目包含前端、后端和 Python 微服务模块，以下为各模块的运行方法。

---

## 目录结构
```

.
 ├── frontend/             # 前端项目
 ├── backend/              # 后端服务
 └── py_mini_service/      # Python 微服务
```
---

## 1. 前端运行

### 进入前端目录

```bash
cd ./frontend
```

### 安装依赖

```bash
yarn install
```

### 启动开发服务器

```bash
yarn dev
```

### 打包构建

```bash
yarn build
```

------

## 2. 后端运行

### 进入后端目录

```bash
cd ./backend
```

### 安装依赖

```bash
go get -u all
```

### 启动服务

```bash
go run main.go
```

------

## 3. Python 微服务运行

### 进入服务目录

```bash
cd ./py_mini_service
```

### 创建并激活虚拟环境（建议使用 conda）

```bash
conda create -n my_env python==3.11
conda activate my_env
```

### 安装依赖

```bash
pip install -r requirements.txt
```

### 拉取基础镜像（如使用 Docker）

```bash
docker pull python:3.12-slim
```

### 启动各 Python 脚本服务

```bash
python code_sandbox/run_service.py
python file_parse/parse.py
python rag/rag.py
```

------

## 注意事项

- 请确保本地已安装好 Node.js、Go、Conda（或 Python3）、Docker 等环境。