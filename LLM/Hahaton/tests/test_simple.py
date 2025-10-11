import requests

print("=== Тест без LLM ===")
payload = {
    "user_query": "Простой тест",
    "source_config": {"type": "csv"},
    "target_config": {"type": "postgresql"},
    "operation_type": "analyze"
}
response = requests.post("http://localhost:8124/api/v1/process", json=payload)
print(f"Status: {response.status_code}")
print(f"Response: {response.text}")
