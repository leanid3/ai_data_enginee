"""Reusable validation helpers."""

from __future__ import annotations

import re
from typing import Iterable


def ensure_non_empty(value: str, field_name: str) -> str:
    if not value:
        raise ValueError(f"{field_name} must not be empty")
    return value


def validate_cron_expression(expression: str) -> str:
    pattern = re.compile(r"^([*\d/,-]+\s*){5,6}$")
    if not pattern.match(expression.strip()):
        raise ValueError("Invalid cron expression")
    return expression


def ensure_unique(items: Iterable[str], field_name: str) -> None:
    seen = set()
    for item in items:
        if item in seen:
            raise ValueError(f"Duplicate value '{item}' in {field_name}")
        seen.add(item)


