import sys
import time
from pathlib import Path
from typing import Any, Dict, List

import asyncio
import httpx
import pytest
from pydantic import BaseModel, Field


PROJECT_ROOT = Path(__file__).resolve().parents[1]
if str(PROJECT_ROOT) not in sys.path:
    sys.path.insert(0, str(PROJECT_ROOT))

from LLM.core.tools import load_llm_tools


class TestResult(BaseModel):
    test_name: str
    success: bool
    response_time: float
    agents_used: List[str]
    confidence_score: float
    errors: List[str]
    response_data: Dict[str, Any]
    tools_reported: List[str] = Field(default_factory=list)
    tools_missing: List[str] = Field(default_factory=list)
    tools_expected: List[str] = Field(default_factory=list)


class DataEngineerTester:
    def __init__(self, base_url: str = "http://localhost:8124"):
        self.base_url = base_url
        self.client = httpx.AsyncClient(timeout=60.0)
        self.results: List[TestResult] = []

    async def test_full_process_request(self) -> TestResult:
        """Тест 1: FULL_PROCESS - сложный CSV с автоматическим обновлением"""
        test_data = {
            "user_query": "У меня есть CSV файл с продажами за последние 3 года, 5 миллионов записей. Нужно загрузить данные в систему для ежедневной аналитики с дашбордами. Файл содержит дату продажи, товар, количество, цену, регион, менеджера. Создай полное решение с автоматическим обновлением каждую ночь.",
            "source_config": {
                "type": "csv",
                "file_path": "sales_3_years.csv",
                "connection_string": None
            },
            "target_config": {
                "type": "clickhouse",
                "database_name": "analytics"
            },
            "data_sample": {
                "headers": ["sale_date", "product", "quantity", "price", "region", "manager"],
                "sample_rows": [
                    {"sale_date": "2024-01-15", "product": "laptop", "quantity": 2, "price": 1500.00, "region": "moscow", "manager": "ivanov"},
                    {"sale_date": "2024-01-16", "product": "phone", "quantity": 1, "price": 800.00, "region": "spb", "manager": "petrov"},
                    {"sale_date": "2024-01-17", "product": "tablet", "quantity": 3, "price": 600.00, "region": "kazan", "manager": "sidorov"}
                ],
                "total_rows": 5000000,
                "file_size": 2147483648
            },
            "operation_type": "full_process"
        }

        return await self._make_request("FULL_PROCESS Test", test_data)

    async def test_analyze_request(self) -> TestResult:
        """Тест 2: ANALYZE - JSON с событиями пользователей"""
        test_data = {
            "user_query": "Проанализируй структуру данных: у меня JSON с событиями пользователей, примерно 100К записей в день. Поля: timestamp, user_id, event_type, properties (вложенный объект), session_id. Определи качество данных и возможные проблемы.",
            "source_config": {
                "type": "json",
                "file_path": "user_events.json"
            },
            "target_config": {
                "type": "auto"
            },
            "data_sample": {
                "headers": ["timestamp", "user_id", "event_type", "properties", "session_id"],
                "sample_rows": [
                    {
                        "timestamp": "2024-01-15T10:30:00Z",
                        "user_id": "12345",
                        "event_type": "page_view",
                        "properties": {"page": "/home", "referrer": "google.com", "duration": 45},
                        "session_id": "abc123"
                    },
                    {
                        "timestamp": "2024-01-15T10:31:15Z",
                        "user_id": "12345",
                        "event_type": "click",
                        "properties": {"element": "button", "text": "Buy Now", "position": {"x": 150, "y": 200}},
                        "session_id": "abc123"
                    }
                ],
                "total_rows": 100000,
                "file_size": 52428800
            },
            "operation_type": "analyze"
        }

        return await self._make_request("ANALYZE Test", test_data)

    async def test_recommend_request(self) -> TestResult:
        """Тест 3: RECOMMEND - выбор БД для логов веб-сервера"""
        test_data = {
            "user_query": "Какую БД выбрать для хранения логов веб-сервера? Примерно 10GB в день, нужна быстрая агрегация по часам/дням, редко обращаемся к данным старше месяца. Важна скорость записи и сжатие.",
            "source_config": {
                "type": "logs",
                "file_path": "access.log"
            },
            "target_config": {
                "type": "auto"
            },
            "data_sample": {
                "headers": ["timestamp", "ip", "method", "url", "status", "size", "user_agent"],
                "sample_rows": [
                    {"timestamp": "2024-01-15T10:30:00Z", "ip": "192.168.1.1", "method": "GET", "url": "/api/users", "status": 200, "size": 1024, "user_agent": "Mozilla/5.0"},
                    {"timestamp": "2024-01-15T10:30:01Z", "ip": "10.0.0.5", "method": "POST", "url": "/api/auth", "status": 401, "size": 256, "user_agent": "curl/7.68.0"}
                ],
                "total_rows": 10000000,
                "file_size": 10737418240
            },
            "metadata": {
                "daily_volume": "10GB",
                "retention": "1 month active, 1 year archive",
                "access_patterns": ["hourly aggregations", "daily reports", "rare historical queries"]
            },
            "operation_type": "recommend"
        }

        return await self._make_request("RECOMMEND Test", test_data)

    async def test_generate_ddl_request(self) -> TestResult:
        """Тест 4: GENERATE_DDL - схема для интернет-магазина"""
        test_data = {
            "user_query": "Создай схему БД для интернет-магазина. Нужны таблицы: товары (10К позиций), заказы (1М в месяц), клиенты (500К), платежи. Данные будут использоваться и для операций, и для аналитики.",
            "source_config": {
                "type": "schema_design"
            },
            "target_config": {
                "type": "postgresql",
                "database_name": "ecommerce"
            },
            "metadata": {
                "tables": {
                    "products": {"count": 10000, "growth": "5% monthly"},
                    "orders": {"count": 1000000, "growth": "monthly"},
                    "customers": {"count": 500000, "growth": "10% monthly"},
                    "payments": {"count": 1200000, "growth": "monthly"}
                },
                "workload": ["OLTP operations", "analytics queries", "reporting"]
            },
            "operation_type": "generate_ddl"
        }

        return await self._make_request("GENERATE_DDL Test", test_data)

    async def test_build_pipeline_request(self) -> TestResult:
        """Тест 5: BUILD_PIPELINE - ETL из PostgreSQL в ClickHouse"""
        test_data = {
            "user_query": "Построй пайплайн для ежечасной выгрузки данных из PostgreSQL (таблица transactions) в ClickHouse с агрегацией по часам. Нужно считать сумму, количество и среднее по типам транзакций.",
            "source_config": {
                "type": "postgresql",
                "connection_string": "postgresql://user:pass@localhost:5432/finance",
                "table_name": "transactions"
            },
            "target_config": {
                "type": "clickhouse",
                "connection_string": "clickhouse://default:@localhost:8123/analytics",
                "table_name": "transaction_hourly_agg"
            },
            "metadata": {
                "schedule": "0 */1 * * *",  # каждый час
                "aggregations": ["SUM(amount)", "COUNT(*)", "AVG(amount)"],
                "group_by": ["transaction_type", "hour"]
            },
            "operation_type": "build_pipeline"
        }

        return await self._make_request("BUILD_PIPELINE Test", test_data)

    async def test_optimize_request(self) -> TestResult:
        """Тест 6: OPTIMIZE - оптимизация сложного запроса"""
        test_data = {
            "user_query": "Оптимизируй запрос: SELECT * FROM orders o JOIN customers c ON o.customer_id = c.id WHERE o.created_at > '2024-01-01' AND c.country = 'Russia' ORDER BY o.total_amount DESC. Таблица orders - 10М записей, customers - 1М.",
            "source_config": {
                "type": "postgresql"
            },
            "target_config": {
                "type": "postgresql"
            },
            "metadata": {
                "query": "SELECT * FROM orders o JOIN customers c ON o.customer_id = c.id WHERE o.created_at > '2024-01-01' AND c.country = 'Russia' ORDER BY o.total_amount DESC",
                "table_stats": {
                    "orders": {"rows": 10000000, "size": "5GB"},
                    "customers": {"rows": 1000000, "size": "500MB"}
                }
            },
            "operation_type": "optimize"
        }

        return await self._make_request("OPTIMIZE Test", test_data)

    async def _make_request(self, test_name: str, data: Dict[str, Any]) -> TestResult:
        start_time = time.time()

        try:
            response = await self.client.post(
                f"{self.base_url}/api/v1/process",
                json=data,
                headers={"Content-Type": "application/json"}
            )
            response_time = time.time() - start_time

            if response.status_code == 200:
                result_data = response.json()

                expected_tools = sorted({tool.name for tool in load_llm_tools()})
                reported_tools = sorted(
                    result_data.get("tools_used")
                    or result_data.get("tools_available")
                    or []
                )

                tool_errors: List[str] = []
                if not reported_tools:
                    tool_errors.append("Response does not provide tool usage information")
                else:
                    unexpected = [tool for tool in reported_tools if tool not in expected_tools]
                    if unexpected:
                        tool_errors.append(
                            "Unexpected tools reported: " + ", ".join(unexpected)
                        )

                missing_tools = [tool for tool in expected_tools if tool not in reported_tools]
                if missing_tools:
                    tool_errors.append(
                        "Missing tools in response: " + ", ".join(missing_tools)
                    )

                combined_errors = result_data.get("errors", []) + tool_errors
                success = not combined_errors

                return TestResult(
                    test_name=test_name,
                    success=success,
                    response_time=response_time,
                    agents_used=result_data.get("agents_used", []),
                    confidence_score=result_data.get("confidence_score", 0.0) if success else 0.0,
                    errors=combined_errors,
                    response_data=result_data,
                    tools_reported=reported_tools,
                    tools_missing=missing_tools,
                    tools_expected=expected_tools,
                )
            else:
                return TestResult(
                    test_name=test_name,
                    success=False,
                    response_time=response_time,
                    agents_used=[],
                    confidence_score=0.0,
                    errors=[f"HTTP {response.status_code}: {response.text}"],
                    response_data={},
                    tools_reported=[],
                    tools_missing=[],
                    tools_expected=[],
                )

        except Exception as e:
            response_time = time.time() - start_time
            return TestResult(
                test_name=test_name,
                success=False,
                response_time=response_time,
                agents_used=[],
                confidence_score=0.0,
                errors=[str(e)],
                response_data={},
                tools_reported=[],
                tools_missing=[],
                tools_expected=[],
            )

    async def test_edge_cases(self) -> List[TestResult]:
        """Тестирование граничных случаев"""
        edge_cases = []

        # Минимальный запрос
        minimal_data = {
            "user_query": "Помоги с данными",
            "source_config": {"type": "unknown"},
            "target_config": {"type": "auto"},
            "operation_type": "analyze"
        }
        edge_cases.append(await self._make_request("Minimal Request", minimal_data))

        # Максимально сложный запрос
        complex_data = {
            "user_query": "Интегрируй данные из 10 источников: 3 PostgreSQL БД, 2 API REST, Kafka топик, S3 bucket с Parquet файлами, ClickHouse, MongoDB и Excel файлы с FTP. Нужна дедупликация, валидация, обогащение справочниками, расчет 50 метрик, инкрементальная загрузка, версионирование, аудит, алерты при аномалиях, автоматическое масштабирование и откат при ошибках.",
            "source_config": {"type": "multi_source"},
            "target_config": {"type": "data_warehouse"},
            "operation_type": "full_process",
            "metadata": {
                "sources": ["postgresql", "api", "kafka", "s3", "clickhouse", "mongodb", "excel"],
                "complexity": "maximum",
                "requirements": ["deduplication", "validation", "enrichment", "metrics", "incremental", "versioning", "audit", "alerts", "scaling", "rollback"]
            }
        }
        edge_cases.append(await self._make_request("Maximum Complexity", complex_data))

        return edge_cases

    async def run_all_tests(self) -> None:
        print("Запуск тестирования AI DATA Engineer...")
        print("=" * 60)

        tests = [
            self.test_full_process_request(),
            self.test_analyze_request(),
            self.test_recommend_request(),
            self.test_generate_ddl_request(),
            self.test_build_pipeline_request(),
            self.test_optimize_request()
        ]

        for test_coro in tests:
            result = await test_coro
            self.results.append(result)
            self._print_test_result(result)

        print("\nТестирование граничных случаев...")
        print("-" * 40)
        edge_results = await self.test_edge_cases()
        self.results.extend(edge_results)

        for result in edge_results:
            self._print_test_result(result)

        self._print_summary()

        await self.client.aclose()

    def _print_test_result(self, result: TestResult) -> None:
        """Вывести результат теста"""
        status = "[PASS]" if result.success else "[FAIL]"
        print(f"{status} {result.test_name}")
        print(f"   Время: {result.response_time:.2f}s")

        if result.success:
            print(f"   Агенты: {', '.join(result.agents_used)}")
            print(f"   Уверенность: {result.confidence_score:.2f}")
            if result.response_data.get("user_report"):
                report = result.response_data["user_report"][:100] + "..." if len(result.response_data["user_report"]) > 100 else result.response_data["user_report"]
                print(f"   Отчет: {report}")
            if result.tools_reported:
                print(f"   Инструменты: {', '.join(result.tools_reported)}")
        else:
            print(f"   Ошибки: {'; '.join(result.errors)}")
            if result.tools_missing:
                print(f"   Пропущенные инструменты: {', '.join(result.tools_missing)}")
        print()

    def _print_summary(self) -> None:
        """Вывести сводку по тестам"""
        print("\nСВОДКА ТЕСТИРОВАНИЯ")
        print("=" * 60)

        total = len(self.results)
        passed = sum(1 for r in self.results if r.success)
        failed = total - passed

        avg_time = sum(r.response_time for r in self.results) / total if total > 0 else 0
        avg_confidence = sum(r.confidence_score for r in self.results if r.success) / passed if passed > 0 else 0

        print(f"Всего тестов: {total}")
        print(f"Успешных: {passed} ({passed/total*100:.1f}%)")
        print(f"Неудачных: {failed} ({failed/total*100:.1f}%)")
        print(f"Среднее время: {avg_time:.2f}s")
        print(f"Средняя уверенность: {avg_confidence:.2f}")

        all_agents = set()
        reported_tools = set()
        missing_tools = set()
        for result in self.results:
            all_agents.update(result.agents_used)
            reported_tools.update(result.tools_reported)
            missing_tools.update(result.tools_missing)

        print(f"\nИспользованные агенты: {', '.join(sorted(all_agents))}")
        print(f"Задекларированные инструменты: {', '.join(sorted(reported_tools))}")
        if missing_tools:
            print(f"Отсутствующие инструменты: {', '.join(sorted(missing_tools))}")

        if failed > 0:
            print(f"\nНеудачные тесты:")
            for result in self.results:
                if not result.success:
                    print(f"   - {result.test_name}: {'; '.join(result.errors)}")


# Pytest тесты
@pytest.fixture
async def tester():
    tester = DataEngineerTester()
    yield tester
    await tester.client.aclose()


@pytest.mark.asyncio
async def test_full_process(tester):
    result = await tester.test_full_process_request()
    assert result.success, f"Test failed: {result.errors}"


@pytest.mark.asyncio
async def test_analyze(tester):
    result = await tester.test_analyze_request()
    assert result.success, f"Test failed: {result.errors}"


@pytest.mark.asyncio
async def test_recommend(tester):
    result = await tester.test_recommend_request()
    assert result.success, f"Test failed: {result.errors}"


@pytest.mark.asyncio
async def test_generate_ddl(tester):
    result = await tester.test_generate_ddl_request()
    assert result.success, f"Test failed: {result.errors}"


@pytest.mark.asyncio
async def test_build_pipeline(tester):
    result = await tester.test_build_pipeline_request()
    assert result.success, f"Test failed: {result.errors}"


@pytest.mark.asyncio
async def test_optimize(tester):
    result = await tester.test_optimize_request()
    assert result.success, f"Test failed: {result.errors}"


@pytest.mark.asyncio
async def test_health_check():
    """Тест здоровья сервиса"""
    async with httpx.AsyncClient(timeout=10.0) as client:
        response = await client.get("http://localhost:8124/api/v1/health")
        assert response.status_code == 200
        data = response.json()
        assert data["status"] == "healthy"


@pytest.mark.asyncio
async def test_agents_list():
    """Тест получения списка агентов"""
    async with httpx.AsyncClient(timeout=10.0) as client:
        response = await client.get("http://localhost:8124/api/v1/agents")
        assert response.status_code == 200
        data = response.json()
        assert "agents" in data
        assert isinstance(data["agents"], list)


if __name__ == "__main__":
    # Запуск как скрипта
    async def main():
        tester = DataEngineerTester()
        try:
            await tester.run_all_tests()
        finally:
            await tester.client.aclose()

    asyncio.run(main())