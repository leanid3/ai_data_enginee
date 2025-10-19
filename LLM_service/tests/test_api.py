import requests
import json

print("=== Тест Health Endpoint ===")
response = requests.get("http://localhost:8124/api/v1/health")
print(f"Status: {response.status_code}")
print(f"Response: {response.json()}")
print()

print("=== Тест загрузки файла ===")
with open("test_data.csv", "rb") as f:
    files = {"file": ("test_data.csv", f, "text/csv")}
    response = requests.post("http://localhost:8124/api/v1/analyze-file", files=files)
    print(f"Status: {response.status_code}")
    if response.status_code == 200:
        print("Успешно!")
        data = response.json()
        print(f"Pipeline ID: {data.get('pipeline_id')}")
        print(f"Success: {data.get('success')}")
    else:
        print(f"Error: {response.text}")
print()

print("=== Тест простого запроса ===")
payload = {
    "user_query": "Создать простой ETL пайплайн",
    "source_config": {"type": "csv", "file_path": "test.csv"},
    "target_config": {"type": "postgresql", "table_name": "test_table"},
    "operation_type": "full_process"
}
response = requests.post("http://localhost:8124/api/v1/process", json=payload)
print(f"Status: {response.status_code}")
if response.status_code == 200:
    print("Успешно!")
    data = response.json()
    print(f"Pipeline ID: {data.get('pipeline_id')}")
    print(f"Success: {data.get('success')}")
else:
    print(f"Error: {response.text}")
