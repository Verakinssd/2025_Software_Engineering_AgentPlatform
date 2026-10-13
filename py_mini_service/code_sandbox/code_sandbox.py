import asyncio
from typing import Optional, List, Dict, Any
import os
import tempfile
from sandbox.client import SANDBOX_CLIENT, SandboxSettings

async def initialize_sandbox(
    use_network: bool = False,
    memory_limit: str = "512m",
    cpu_limit: float = 1.0,
    timeout: int = 300,
    custom_image: Optional[str] = None,
) -> str:
    """初始化一个代码执行沙盒环境
    
    Args:
        use_network: 是否允许网络访问，默认为False
        memory_limit: 内存限制，默认为512m
        cpu_limit: CPU限制，默认为1.0核心
        timeout: 默认命令超时时间（秒），默认为300秒
        custom_image: 自定义Docker镜像，默认为None（使用python:3.12-slim） docker pull python:3.12-slim
        
    Returns:
        str: 初始化结果信息
    """
    config = SandboxSettings(
        use_sandbox=True,
        network_enabled=use_network,
        memory_limit=memory_limit,
        cpu_limit=cpu_limit,
        timeout=timeout,
    )
    
    if custom_image:
        config.image = custom_image
    
    try:
        await SANDBOX_CLIENT.create(config=config)
        return "沙盒环境初始化成功"
    except Exception as e:
        return f"沙盒环境初始化失败: {str(e)}"


async def cleanup_sandbox() -> str:
    """清理并关闭沙盒环境
    
    Returns:
        str: 清理结果信息
    """
    try:
        await SANDBOX_CLIENT.cleanup()
        return "沙盒环境已清理"
    except Exception as e:
        return f"沙盒环境清理失败: {str(e)}"


async def run_code(code: str, timeout: Optional[int] = None) -> str:
    """在沙盒环境中执行Python代码
    
    Args:
        code: 要执行的Python代码
        timeout: 执行超时时间（秒），不指定则使用沙盒默认值
        
    Returns:
        str: 代码执行的输出结果
    """
    try:
        # 将代码写入临时文件
        await SANDBOX_CLIENT.write_file("/workspace/temp_script.py", code)
        
        # 执行代码
        result = await SANDBOX_CLIENT.run_command("python /workspace/temp_script.py", timeout)
        return result
    except Exception as e:
        return f"代码执行失败: {str(e)}"


async def execute_shell_command(command: str, timeout: Optional[int] = None) -> str:
    """在沙盒环境中执行Shell命令
    
    Args:
        command: 要执行的Shell命令
        timeout: 执行超时时间（秒），不指定则使用沙盒默认值
        
    Returns:
        str: 命令执行的输出结果
    """
    try:
        result = await SANDBOX_CLIENT.run_command(command, timeout)
        return result
    except Exception as e:
        return f"命令执行失败: {str(e)}"


async def write_file_to_sandbox(file_path: str, content: str) -> str:
    """向沙盒环境中写入文件
    
    Args:
        file_path: 沙盒中的文件路径
        content: 文件内容
        
    Returns:
        str: 操作结果
    """
    try:
        await SANDBOX_CLIENT.write_file(file_path, content)
        return f"文件已成功写入到 {file_path}"
    except Exception as e:
        return f"文件写入失败: {str(e)}"


async def read_file_from_sandbox(file_path: str) -> str:
    """从沙盒环境中读取文件
    
    Args:
        file_path: 沙盒中的文件路径
        
    Returns:
        str: 文件内容
    """
    try:
        content = await SANDBOX_CLIENT.read_file(file_path)
        return content
    except Exception as e:
        return f"文件读取失败: {str(e)}"


async def run_data_analysis_code(code: str, data_file_path: Optional[str] = None, data_content: Optional[str] = None) -> str:
    """在沙盒环境中执行数据分析代码
    
    Args:
        code: 要执行的Python数据分析代码
        data_file_path: 数据文件路径（可选），如果提供数据内容，则忽略此参数
        data_content: 数据内容（可选），将写入到沙盒中的/workspace/data.csv
        
    Returns:
        str: 数据分析结果
    """
    try:
        await SANDBOX_CLIENT.run_command("pip install pandas numpy matplotlib seaborn scikit-learn")
        
        if data_content:
            await SANDBOX_CLIENT.write_file("/workspace/data.csv", data_content)
        elif data_file_path:
            if not os.path.exists(data_file_path):
                return f"数据文件 {data_file_path} 不存在"
            
            with open(data_file_path, 'r') as f:
                data_content = f.read()
                await SANDBOX_CLIENT.write_file("/workspace/data.csv", data_content)
        
        await SANDBOX_CLIENT.write_file("/workspace/analysis.py", code)
        result = await SANDBOX_CLIENT.run_command("python /workspace/analysis.py")
        return result
    except Exception as e:
        return f"数据分析执行失败: {str(e)}"
    

def initialize_sandbox_sync(
    use_network: bool = False,
    memory_limit: str = "512m", 
    cpu_limit: float = 1.0,
    timeout: int = 300,
    custom_image: Optional[str] = None
) -> str:
    """同步版本的初始化沙盒环境"""
    return asyncio.run(initialize_sandbox(
        use_network=use_network,
        memory_limit=memory_limit,
        cpu_limit=cpu_limit,
        timeout=timeout,
        custom_image=custom_image
    ))


def cleanup_sandbox_sync() -> str:
    """同步版本的清理沙盒环境"""
    return asyncio.run(cleanup_sandbox())


def run_code_sync(
    code: str,
    command: Optional[str] = None,
    files: Optional[Dict[str, str]] = None,
    timeout: Optional[int] = None
) -> str:
    """同步版本的代码执行"""
    return asyncio.run(run_code(
        code=code,
        command=command,
        files=files,
        timeout=timeout
    ))


if __name__ == "__main__":
    # 测试示例
    # 初始化沙盒环境
    print("正在初始化沙盒环境...")
    initialize_sandbox_sync()

    # 测试简单的Python代码执行
    test_code = """
print('Hello from sandbox!')
x = 1 + 2
print(f'计算结果: {x}')
"""
    print("\n执行简单Python代码:")
    result = run_code_sync(test_code)
    print(result)

    # 测试文件操作
    test_code_with_file = """
with open('/workspace/test.txt', 'w') as f:
    f.write('这是测试文件内容')
    
with open('/workspace/test.txt', 'r') as f:
    content = f.read()
print(f'读取文件内容: {content}')
"""
    print("\n执行文件操作测试:")
    result = run_code_sync(test_code_with_file)
    print(result)

    # 清理沙盒环境
    print("\n清理沙盒环境...")
    cleanup_sandbox_sync()
