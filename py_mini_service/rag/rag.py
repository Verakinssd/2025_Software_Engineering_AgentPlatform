import os
from typing import List
from rank_bm25 import BM25Okapi
from langchain_text_splitters import RecursiveCharacterTextSplitter
from openai import OpenAI
import jieba
from fastapi import FastAPI, Form
import uvicorn

app = FastAPI()

def split_doc(text: str,
              chunk_size: int = 512,
              chunk_overlap: int = 64) -> List[str]:
    splitter = RecursiveCharacterTextSplitter(
        chunk_size=chunk_size,
        chunk_overlap=chunk_overlap,
        separators=["\n\n", "\n", "。", ".", " ", ""]
    )
    return splitter.split_text(text)

def build_bm25(chunks: List[str]) -> BM25Okapi:
    tokenized = [list(jieba.cut(c)) for c in chunks]
    return BM25Okapi(tokenized)

def retrieve(query: str,
             bm25: BM25Okapi,
             chunks: List[str],
             k: int = 3) -> List[str]:
    tokenized_q = list(jieba.cut(query))
    scores = bm25.get_scores(tokenized_q)
    topk_idx = sorted(range(len(scores)), key=lambda i: scores[i], reverse=True)[:k]
    return [chunks[i] for i in topk_idx]

def generate_answer(query: str,
                    context_chunks: List[str],
                    client: OpenAI,
                    model: str = "glm-4-flash-250414") -> str:
    context = "\n\n".join(f"[{i+1}] {c}" for i, c in enumerate(context_chunks))
    messages = [
        {"role": "user",
         "content": f"问题：{query}\n\n参考资料：\n{context}"
        }
    ]
    # print(context)
    # print("\n\n--------------------------------\n\n")
    resp = client.chat.completions.create(
        model=model,
        messages=messages,
        temperature=0.1,
        max_tokens=512
    )
    return resp.choices[0].message.content.strip()

def rag_answer(document: str, query: str) -> str:
    chunks = split_doc(document)
    bm25   = build_bm25(chunks)
    top    = retrieve(query, bm25, chunks, k=3)

    client = OpenAI(
        api_key=os.getenv("ZHIPUAI_API_KEY"),
        base_url="https://open.bigmodel.cn/api/paas/v4/"
    )
    return generate_answer(query, top, client)

@app.post("/api/query/")
async def query(query: str = Form(...), doc: str = Form(...)):
    try:
        res = rag_answer(document=doc, query=query)
        return {"answer": res}
    except Exception as e:
        return {"answer": f"发生错误：{str(e)}"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=5002)
    # document = open("../file_parse/test_files/1.txt", "r", encoding="utf-8").read()
    # query    = "问题 7-15 的答案是？"
    # answer   = rag_answer(document, query)
    # print(answer)
