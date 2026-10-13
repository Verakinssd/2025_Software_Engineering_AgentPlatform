from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.responses import JSONResponse
import uvicorn
import os
import shutil
from markitdown import MarkItDown

app = FastAPI()

# 实例化 MarkItDown 转换器
md = MarkItDown()

def parse_file_to_text(file_path: str) -> str:
    """
    使用 MarkItDown 将给定文件解析为文本。
    """
    result = md.convert(file_path)
    return result.text_content

@app.post("/api/parse/")
async def parse_doc(file: UploadFile = File(...)):
    # 临时保存上传的文件
    temp_dir = "temp_uploads"
    os.makedirs(temp_dir, exist_ok=True)
    temp_path = os.path.join(temp_dir, file.filename)
    try:
        with open(temp_path, "wb") as buffer:
            shutil.copyfileobj(file.file, buffer)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"保存上传文件失败: {e}")

    # 文件大小
    size = os.path.getsize(temp_path)

    try:
        text = parse_file_to_text(temp_path)
    except Exception as e:
        # 解析失败，返回 400
        raise HTTPException(status_code=400, detail=f"解析文件失败: {e}")
    finally:
        # 清理临时文件
        try:
            os.remove(temp_path)
        except OSError:
            pass

    # 构建响应
    response = {
        "filename": file.filename,
        "size": size,
        "text": text
    }
    return JSONResponse(content=response)

if __name__ == "__main__":
    # print(f"""
    # docs_url: {app.docs_url}
    # redoc_url: {app.redoc_url}
    # openapi_url: {app.openapi_url}
    # Registered routes:
    # """)
    # for r in app.routes:
    #     print(f"  {r.methods}  -> {r.path}")
    uvicorn.run(app, host="0.0.0.0", port=5001)

    # text = parse_file_to_text(r"test_files\1.doc")
    # print(text)
