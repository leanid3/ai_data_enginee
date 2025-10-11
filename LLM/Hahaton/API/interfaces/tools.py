from LLM.core import Tool
import asyncio

async def mock_web_search(query: str) -> str:
    """Имитация веб-поиска"""
    return f"Результаты поиска по запросу '{query}': [Пример 1], [Пример 2]"

async def mock_sftp_listdir(path: str) -> list:
    """Имитация списка файлов на SFTP"""
    return ["data.csv", "config.json", "logs.txt"]

async def mock_run_python_code(code: str) -> str:
    """Имитация выполнения Python-кода"""
    return f"Код выполнен. Вывод: 42"

def get_default_tools() -> list[Tool]:
    return [
        Tool(
            func=mock_web_search,
            name="web_search",
            description="Поиск в интернете по заданному запросу",
            input_schema={
                "query": {
                    "type": "string",
                    "description": "Поисковый запрос на русском или английском языке"
                }
            }
        ),
        Tool(
            func=mock_sftp_listdir,
            name="sftp_listdir",
            description="Получить список файлов в указанной директории на SFTP-сервере",
            input_schema={
                "path": {
                    "type": "string",
                    "description": "Путь к директории на SFTP (например, /data/)"
                }
            }
        ),
        Tool(
            func=mock_run_python_code,
            name="python_interpreter",
            description="Выполнить переданный Python-код и вернуть результат",
            input_schema={
                "code": {
                    "type": "string",
                    "description": "Корректный Python-код для выполнения"
                }
            }
        ),
    ]