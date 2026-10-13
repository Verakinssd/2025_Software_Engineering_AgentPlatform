from fastapi import FastAPI, HTTPException, Body
from pydantic import BaseModel
from typing import Optional
import asyncio
from code_sandbox import run_code, initialize_sandbox, cleanup_sandbox

app = FastAPI(title="代码沙盒服务", description="提供Python代码执行的安全沙盒环境")

class CodeExecutionRequest(BaseModel):
    """代码执行请求模型"""
    code: str
    timeout: Optional[int] = None

class CodeExecutionResponse(BaseModel):
    """代码执行响应模型"""
    result: str
    success: bool

@app.on_event("startup")
async def startup_event():
    """服务启动时初始化沙盒环境"""
    await initialize_sandbox(
        use_network=False,
        memory_limit="512m",
        cpu_limit=1.0,
        timeout=300
    )

@app.on_event("shutdown")
async def shutdown_event():
    """服务关闭时清理沙盒环境"""
    await cleanup_sandbox()

@app.post("/execute", response_model=CodeExecutionResponse)
async def execute_code(request: CodeExecutionRequest = Body(...)):
    """
    执行Python代码并返回结果
    
    Args:
        request: 包含代码和可选的超时设置
        
    Returns:
        执行结果
    """
    try:
        result = await run_code(request.code, request.timeout)
        return CodeExecutionResponse(result=result, success=True)
    except Exception as e:
        return CodeExecutionResponse(result=f"执行错误: {str(e)}", success=False)

@app.get("/health")
async def health_check():
    """健康检查接口"""
    return {"status": "healthy"}

# 用于直接运行服务
if __name__ == "__main__":
    import uvicorn
    uvicorn.run("code_sandbox_service:app", host="0.0.0.0", port=5000, reload=True)
