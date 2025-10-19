"""Tool for formatting analysis results into Markdown reports."""

from __future__ import annotations

from typing import Dict, List

from markdown import markdown
from tabulate import tabulate

from ..base import BaseTool, ToolInput, ToolOutput


class ReportSection(ToolInput):
    name: str
    content: Dict | List | str


class FormatMarkdownReportInput(ToolInput):
    title: str
    sections: List[ReportSection]
    include_toc: bool = True
    include_mermaid: bool = True
    include_code: bool = True
    language: str = "ru"


class FormatMarkdownReportOutput(ToolOutput):
    markdown: str
    html: str | None = None
    pdf_path: str | None = None
    word_count: int = 0
    sections_count: int = 0
    has_diagrams: bool = False
    has_code: bool = False


class FormatMarkdownReportTool(
    BaseTool[FormatMarkdownReportInput, FormatMarkdownReportOutput]
):
    name = "format_markdown_report"
    description = "Format structured results into a Markdown report"
    input_schema = FormatMarkdownReportInput
    output_schema = FormatMarkdownReportOutput

    async def execute(self, params: FormatMarkdownReportInput) -> FormatMarkdownReportOutput:
        markdown_content = self._build_markdown(params)
        html_content = markdown(markdown_content)
        word_count = len(markdown_content.split())

        return FormatMarkdownReportOutput(
            success=True,
            markdown=markdown_content,
            html=html_content,
            word_count=word_count,
            sections_count=len(params.sections),
            has_diagrams=params.include_mermaid,
            has_code=params.include_code,
        )

    def _build_markdown(self, params: FormatMarkdownReportInput) -> str:
        parts = [f"# {params.title}"]
        if params.include_toc:
            parts.append("## Содержание")
            for section in params.sections:
                parts.append(f"- [{section.name}](#{section.name.lower().replace(' ', '-')})")

        for section in params.sections:
            parts.append(f"## {section.name}")
            parts.append(self._render_content(section.content))

        return "\n\n".join(parts)

    def _render_content(self, content) -> str:
        if isinstance(content, str):
            return content
        if isinstance(content, list):
            return "\n".join(f"- {item}" for item in content)
        if isinstance(content, dict):
            table = tabulate(content.items(), tablefmt="github")
            return table
        return str(content)


