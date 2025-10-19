import io
import json
from typing import Any, Dict, List
import pandas as pd
from API.schemas import DataSample


class DataProcessor:
    async def __call__(self, file_type, file_content: bytes, filename: str):
        if file_type == 'csv':
            return await self.analyze_csv(file_content, filename)
        elif file_type == 'json':
            return await self.analyze_json(file_content, filename)
        else:
            raise FileNotFoundError("Unsupported file type")

    @staticmethod
    async def analyze_csv(file_content: bytes, filename: str) -> DataSample:
        df = pd.read_csv(io.BytesIO(file_content))
        return DataSample(
            headers=df.columns.tolist(),
            sample_rows=df.head(10).to_dict("records"),
            total_rows=len(df),
            file_size=len(file_content),
        )

    @staticmethod
    async def analyze_json(file_content: bytes, filename: str) -> DataSample:
        data = json.loads(file_content)
        if isinstance(data, list) and data:
            headers = list(data[0].keys()) if isinstance(data[0], dict) else ["value"]
            sample_rows = data[:10]
            total_rows = len(data)
        elif isinstance(data, dict):
            headers = list(data.keys())
            sample_rows = [data]
            total_rows = 1
        else:
            headers = ["value"]
            sample_rows = [{"value": str(data)}]
            total_rows = 1
        return DataSample(
            headers=headers,
            sample_rows=sample_rows,
            total_rows=total_rows,
            file_size=len(file_content),
        )