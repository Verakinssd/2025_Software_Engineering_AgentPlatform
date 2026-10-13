# 代码沙盒服务

这是一个用于安全执行Python代码的沙盒服务。它使用Docker容器隔离环境，限制资源使用，保证安全执行用户提交的代码。

## 功能特点

- 安全执行Python代码
- 资源限制（CPU、内存、运行时间）
- 网络访问控制
- 提供REST API接口

## 安装要求

- Python 3.8+
- Docker
- FastAPI
- Uvicorn

## 安装

1. 确保已安装Docker并启动Docker服务

2. 安装Python依赖：

```bash
pip install fastapi uvicorn docker pydantic
```

## 使用方法

### 启动服务

```bash
cd py_mini_service/code_sandbox
python run_service.py
```

服务默认在 http://localhost:8000 启动

### API文档

访问 http://localhost:8000/docs 查看完整的API文档

### 执行代码示例

使用HTTP POST请求执行代码：

```python
import requests

response = requests.post(
    "http://localhost:8000/execute",
    json={
        "code": "print('Hello, World!')\nx = 10\nprint(f'x = {x}')",
        "timeout": 30
    }
)

print(response.json())
```

返回结果示例：

```json
{
    "result": "Hello, World!\nx = 10\n",
    "success": true
}
```

## 安全配置

默认情况下：
- 禁用网络访问
- 内存限制: 512MB
- CPU限制: 1核心
- 执行超时: 300秒

## 注意事项

- 该服务仅用于安全执行代码，不应用于生产环境的关键任务
- 沙盒环境中的文件在会话结束后会被清理
- 如需持久化数据，请使用外部存储服务 