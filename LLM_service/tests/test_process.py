def test_process_minimal(client):
    payload = {
        "user_query": "Построить пайплайн",
        "source_config": {"type": "csv"},
        "target_config": {"type": "clickhouse"},
        "operation_type": "full_process",
    }
    r = client.post("/api/v1/process", json=payload)
    assert r.status_code == 200
    data = r.json()
    assert data["success"] is True
    assert "pipeline_id" in data
