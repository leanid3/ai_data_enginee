import os
from dotenv import load_dotenv

print("=== До load_dotenv ===")
print("DEFAULT_MODEL:", os.getenv("DEFAULT_MODEL"))

load_dotenv()

print("=== После load_dotenv ===")
print("DEFAULT_MODEL:", os.getenv("DEFAULT_MODEL"))
print("OPENROUTER_API_KEY:", os.getenv("OPENROUTER_API_KEY"))

print("=== Содержимое .env ===")
try:
    with open(".env", "r", encoding="utf-8") as f:
        print(f.read())
except Exception as e:
    print(f"Ошибка чтения .env: {e}")
