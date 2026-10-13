import uvicorn

if __name__ == "__main__":
    # 启动代码沙盒服务
    uvicorn.run("code_sandbox_service:app", 
                host="0.0.0.0", 
                port=5000, 
                reload=True,
                log_level="info")
    
    print("代码沙盒服务已启动，访问 http://localhost:8000/docs 查看API文档") 