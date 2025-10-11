import os
from dotenv import load_dotenv

load_dotenv()

print("OPENROUTER_API_KEY:", os.getenv("OPENROUTER_API_KEY"))
print("OPENROUTER_BASE_URL:", os.getenv("OPENROUTER_BASE_URL"))
print("DEFAULT_MODEL:", os.getenv("DEFAULT_MODEL"))
