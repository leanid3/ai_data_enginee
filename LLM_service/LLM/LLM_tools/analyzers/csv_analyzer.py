"""CSV结构分析工具实现。
用于分析CSV文件的结构、数据类型和统计信息的高性能工具。
支持自动检测编码、分隔符和模式识别。
"""

from __future__ import annotations

import csv
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Dict, List, Optional

import chardet
import numpy as np
import pandas as pd
from pandas.api.types import is_numeric_dtype

from LLM_service.utils.logger import AutoClassFuncLogger

from pydantic import BaseModel

from ..base import BaseTool, ToolInput, ToolOutput


class AnalyzeCsvInput(ToolInput):
    file_path: str
    sample_size: int = 1000
    encoding: Optional[str] = None
    delimiter: Optional[str] = None


class ColumnStatistics(BaseModel):
    name: str
    data_type: str
    nullable: bool
    unique_count: int
    null_count: int
    null_percentage: float
    sample_values: List[Any]
    statistics: Dict[str, float] | None = None
    pattern: Optional[str] = None
    is_potential_key: bool = False


class AnalyzeCsvOutput(ToolOutput):
    columns: List[ColumnStatistics]
    row_count: int
    file_size_mb: float
    encoding: str
    delimiter: str
    has_header: bool
    data_quality_score: float


# CSV元数据类，存储文件的基本属性
# 包括编码格式、分隔符和表头信息
@dataclass
class CsvMetadata:
    encoding: str  # 文件编码（如UTF-8、GBK等）
    delimiter: str  # 字段分隔符
    has_header: bool  # 是否包含表头


# CSV结构分析工具类
# 提供完整的CSV文件分析功能，包括：
# - 自动检测文件编码和分隔符
# - 分析数据类型和模式
# - 生成数据质量报告
# - 识别潜在的主键列
class AnalyzeCsvStructureTool(BaseTool[AnalyzeCsvInput, AnalyzeCsvOutput]):
    name = "analyze_csv_structure"
    description = "Analyze CSV files to determine their structure and statistics"
    input_schema = AnalyzeCsvInput
    output_schema = AnalyzeCsvOutput

    def __init__(self, logger: AutoClassFuncLogger | None = None) -> None:
        super().__init__(logger)

    async def execute(self, params: AnalyzeCsvInput) -> AnalyzeCsvOutput:
        path = Path(params.file_path)
        if not path.exists():
            raise FileNotFoundError(f"CSV file not found: {path}")

        metadata = self._detect_metadata(path, params)
        df = await self.run_in_thread(
            self._load_dataframe,
            path,
            params.sample_size,
            metadata.encoding,
            metadata.delimiter,
        )

        column_stats = self._analyze_columns(df)
        file_size_mb = path.stat().st_size / (1024 * 1024)
        data_quality_score = self._calculate_quality_score(column_stats)

        return AnalyzeCsvOutput(
            success=True,
            columns=column_stats,
            row_count=len(df),
            file_size_mb=file_size_mb,
            encoding=metadata.encoding,
            delimiter=metadata.delimiter,
            has_header=metadata.has_header,
            data_quality_score=data_quality_score,
        )

    # 检测CSV文件的元数据信息
    # 包括文件编码、分隔符和表头
    # 使用智能检测算法自动识别文件特征
    def _detect_metadata(  # noqa: PLR0913 - clarity
        self,
        path: Path,
        params: AnalyzeCsvInput,
    ) -> CsvMetadata:
        encoding = params.encoding or self._detect_encoding(path)
        delimiter = params.delimiter or self._detect_delimiter(path, encoding)
        has_header = self._has_header(path, delimiter, encoding)
        return CsvMetadata(encoding=encoding, delimiter=delimiter, has_header=has_header)

    def _detect_encoding(self, path: Path) -> str:
        with path.open("rb") as file:
            raw = file.read(64 * 1024)
            result = chardet.detect(raw)
        return result.get("encoding", "utf-8")

    def _detect_delimiter(self, path: Path, encoding: str) -> str:
        with path.open("r", encoding=encoding, errors="ignore") as file:
            dialect = csv.Sniffer().sniff(file.read(2048))
        return dialect.delimiter

    def _has_header(self, path: Path, delimiter: str, encoding: str) -> bool:
        with path.open("r", encoding=encoding, errors="ignore") as file:
            return csv.Sniffer().has_header(file.read(2048))

    def _load_dataframe(
        self,
        path: Path,
        sample_size: int,
        encoding: str,
        delimiter: str,
    ) -> pd.DataFrame:
        return pd.read_csv(
            path,
            delimiter=delimiter,
            encoding=encoding,
            nrows=sample_size,
            low_memory=False,
        )

    # 分析数据列的统计特征
    # 对每一列进行深度分析：
    # - 数据类型推断
    # - 空值统计
    # - 唯一值分析
    # - 数值型数据统计
    # - 模式识别
    def _analyze_columns(self, df: pd.DataFrame) -> List[ColumnStatistics]:
        stats = []
        for column in df.columns:
            series = df[column]
            statistics = None

            if is_numeric_dtype(series):
                statistics = self._numeric_statistics(series)

            stat = ColumnStatistics(
                name=str(column),
                data_type=str(series.dtype),
                nullable=series.isnull().any(),
                unique_count=series.nunique(dropna=True),
                null_count=series.isnull().sum(),
                null_percentage=float(series.isnull().mean() * 100),
                sample_values=self._sample_values(series),
                statistics=statistics,
                pattern=self._detect_pattern(series),
                is_potential_key=self._is_potential_key(series),
            )
            stats.append(stat)
        return stats

    def _numeric_statistics(self, series: pd.Series) -> Dict[str, float]:
        values = series.dropna().astype(float)
        if values.empty:
            return {}
        return {
            "min": float(values.min()),
            "max": float(values.max()),
            "mean": float(values.mean()),
            "median": float(values.median()),
        }

    def _sample_values(self, series: pd.Series, count: int = 5) -> List[Any]:
        sample = series.dropna().head(count).tolist()
        return sample

    def _detect_pattern(self, series: pd.Series) -> Optional[str]:
        sample = series.dropna().astype(str).head(20)
        if sample.empty:
            return None
        if sample.str.fullmatch(r"\d{4}-\d{2}-\d{2}").all():
            return "date"
        if sample.str.contains("@", na=False).all():
            return "email"
        if sample.str.fullmatch(r"\+?\d+[\d\-\s]*").all():
            return "phone"
        if sample.str.fullmatch(r"[A-Za-z0-9_-]+").all():
            return "id"
        return None

    def _is_potential_key(self, series: pd.Series) -> bool:
        return series.is_unique and not series.isnull().any()

    # 计算数据质量得分
    # 基于以下指标：
    # - 数据完整性（空值率）
    # - 数据唯一性
    # - 列类型分布
    def _calculate_quality_score(self, columns: List[ColumnStatistics]) -> float:
        if not columns:
            return 0.0
        weights = []
        for column in columns:
            completeness = 1.0 - column.null_percentage / 100
            uniqueness = min(column.unique_count / max(1, len(column.sample_values)), 1.0)
            weights.append((completeness + uniqueness) / 2)
        return float(np.mean(weights))


