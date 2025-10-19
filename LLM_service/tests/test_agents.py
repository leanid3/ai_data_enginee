"""
Тестовый скрипт для проверки всех ролей агентов системы автоматизации инженерии данных.
"""

import asyncio
import json
import logging
from typing import Dict, Any, List
import sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).parent.parent))

from LLM.core.agents import AgentRole, LLMChainManager, BaseAgent, PromptTemplates
from API.interfaces.clients.open_router_client import OpenRouterClient


logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)


class AgentTester:
    def __init__(self):
        self.llm_client = OpenRouterClient()
        self.chain_manager = LLMChainManager(self.llm_client)
        
    async def test_router_agent(self) -> Dict[str, Any]:
        logger.info("Тестирование ROUTER агента...")
        
        context = {
            "user_request": "Нужно проанализировать CSV файл с продажами и создать хранилище данных",
            "operation_type": "full_process"
        }
        
        try:
            result = await self.chain_manager.call_agent(AgentRole.ROUTER, context)
            logger.info(f"ROUTER результат: {json.dumps(result, indent=2, ensure_ascii=False)}")
            return result
        except Exception as e:
            logger.error(f"Ошибка в ROUTER: {e}")
            return {"error": str(e)}

    async def test_data_analyzer_agent(self) -> Dict[str, Any]:
        logger.info("Тестирование DATA_ANALYZER агента...")
        
        context = {
            "data_sample": {
                "headers": ["order_id", "customer_id", "product_name", "quantity", "price", "order_date"],
                "sample_rows": [
                    {"order_id": 1001, "customer_id": "C001", "product_name": "Laptop", "quantity": 2, "price": 1500.00, "order_date": "2024-01-15"},
                    {"order_id": 1002, "customer_id": "C002", "product_name": "Mouse", "quantity": 10, "price": 25.99, "order_date": "2024-01-16"}
                ],
                "total_rows": 50000,
                "file_size": 5242880
            }
        }
        
        try:
            result = await self.chain_manager.call_agent(AgentRole.DATA_ANALYZER, context)
            logger.info(f"DATA_ANALYZER результат: {json.dumps(result, indent=2, ensure_ascii=False)}")
            return result
        except Exception as e:
            logger.error(f"Ошибка в DATA_ANALYZER: {e}")
            return {"error": str(e)}

    async def test_db_selector_agent(self) -> Dict[str, Any]:
        """Тест селектора БД"""
        logger.info("Тестирование DB_SELECTOR агента...")
        
        context = {
            "analysis_results": {
                "data_type": "OLTP",
                "characteristics": {
                    "volume": "medium",
                    "update_frequency": "daily"
                },
                "structure": {
                    "key_fields": ["order_id", "customer_id"],
                    "partition_fields": ["order_date"]
                },
                "quality_score": 0.85
            }
        }
        
        try:
            result = await self.chain_manager.call_agent(AgentRole.DB_SELECTOR, context)
            logger.info(f"DB_SELECTOR результат: {json.dumps(result, indent=2, ensure_ascii=False)}")
            return result
        except Exception as e:
            logger.error(f"Ошибка в DB_SELECTOR: {e}")
            return {"error": str(e)}

    async def test_ddl_generator_agent(self) -> Dict[str, Any]:
        logger.info("Тестирование DDL_GENERATOR агента...")
        
        context = {
            "target_db": "PostgreSQL",
            "data_structure": {
                "key_fields": ["order_id", "customer_id"],
                "partition_fields": ["order_date"]
            },
            "storage_config": {
                "partitioning": "date_based",
                "replication": 1
            }
        }
        
        try:
            result = await self.chain_manager.call_agent(AgentRole.DDL_GENERATOR, context)
            logger.info(f"DDL_GENERATOR результат: {json.dumps(result, indent=2, ensure_ascii=False)}")
            return result
        except Exception as e:
            logger.error(f"Ошибка в DDL_GENERATOR: {e}")
            return {"error": str(e)}

    async def test_etl_builder_agent(self) -> Dict[str, Any]:
        logger.info("Тестирование ETL_BUILDER агента...")
        
        context = {
            "source_config": {
                "type": "csv",
                "file_path": "/data/sales.csv"
            },
            "target_config": {
                "type": "postgresql",
                "connection_string": "postgresql://user:pass@localhost/sales_db",
                "table_name": "orders"
            },
            "schedule": "daily"
        }
        
        try:
            result = await self.chain_manager.call_agent(AgentRole.ETL_BUILDER, context)
            logger.info(f"ETL_BUILDER результат: {json.dumps(result, indent=2, ensure_ascii=False)}")
            return result
        except Exception as e:
            logger.error(f"Ошибка в ETL_BUILDER: {e}")
            return {"error": str(e)}

    async def test_query_optimizer_agent(self) -> Dict[str, Any]:
        logger.info("Тестирование QUERY_OPTIMIZER агента...")
        
        context = {
            "db_type": "PostgreSQL",
            "queries": [
                "SELECT * FROM orders WHERE order_date BETWEEN '2024-01-01' AND '2024-12-31'",
                "SELECT customer_id, COUNT(*) FROM orders GROUP BY customer_id"
            ],
            "data_volume": "50GB"
        }
        
        try:
            result = await self.chain_manager.call_agent(AgentRole.QUERY_OPTIMIZER, context)
            logger.info(f"QUERY_OPTIMIZER результат: {json.dumps(result, indent=2, ensure_ascii=False)}")
            return result
        except Exception as e:
            logger.error(f"Ошибка в QUERY_OPTIMIZER: {e}")
            return {"error": str(e)}

    async def test_report_generator_agent(self) -> Dict[str, Any]:
        logger.info("Тестирование REPORT_GENERATOR агента...")
        
        context = {
            "all_results": {
                "analysis": {"data_type": "OLTP", "quality_score": 0.85},
                "storage": {"recommended_storage": "PostgreSQL"},
                "ddl": {"ddl_scripts": [{"type": "TABLE", "name": "orders"}]},
                "etl": {"dag_id": "sales_etl", "schedule": "0 2 * * *"},
                "optimization": {"estimated_improvement": "40%"}
            }
        }
        
        try:
            result = await self.chain_manager.call_agent(AgentRole.REPORT_GENERATOR, context)
            logger.info(f"REPORT_GENERATOR результат: {json.dumps(result, indent=2, ensure_ascii=False)}")
            return result
        except Exception as e:
            logger.error(f"Ошибка в REPORT_GENERATOR: {e}")
            return {"error": str(e)}

    async def run_all_tests(self) -> Dict[str, Any]:
        logger.info("Запуск тестирования всех агентов...")
        
        test_results = {}
        
        # Список всех тестов
        test_methods = [
            ("ROUTER", self.test_router_agent),
            ("DATA_ANALYZER", self.test_data_analyzer_agent),
            ("DB_SELECTOR", self.test_db_selector_agent),
            ("DDL_GENERATOR", self.test_ddl_generator_agent),
            ("ETL_BUILDER", self.test_etl_builder_agent),
            ("QUERY_OPTIMIZER", self.test_query_optimizer_agent),
            ("REPORT_GENERATOR", self.test_report_generator_agent)
        ]
        
        # Выполнение тестов
        for agent_name, test_method in test_methods:
            logger.info(f"\n{'='*60}")
            logger.info(f"ТЕСТИРОВАНИЕ АГЕНТА: {agent_name}")
            logger.info(f"{'='*60}")
            
            try:
                result = await test_method()
                test_results[agent_name] = {
                    "status": "success" if "error" not in result else "failed",
                    "result": result
                }
                
                if "error" not in result:
                    logger.info(f"✅ Агент {agent_name} работает корректно")
                else:
                    logger.error(f"❌ Агент {agent_name} завершился с ошибкой")
                    
            except Exception as e:
                logger.error(f"❌ Критическая ошибка при тестировании {agent_name}: {e}")
                test_results[agent_name] = {
                    "status": "critical_error",
                    "error": str(e)
                }
        
        # Сводка результатов
        logger.info(f"\n{'='*60}")
        logger.info("СВОДКА РЕЗУЛЬТАТОВ ТЕСТИРОВАНИЯ")
        logger.info(f"{'='*60}")
        
        successful = sum(1 for result in test_results.values() if result["status"] == "success")
        failed = len(test_results) - successful
        
        logger.info(f"Всего агентов протестировано: {len(test_results)}")
        logger.info(f"Успешно: {successful}")
        logger.info(f"С ошибками: {failed}")
        
        for agent_name, result in test_results.items():
            status_emoji = "✅" if result["status"] == "success" else "❌"
            logger.info(f"{status_emoji} {agent_name}: {result['status']}")
        
        return test_results


async def main():
    print("🚀 Запуск тестирования системы агентов")
    print("=" * 60)
    
    try:
        tester = AgentTester()
        results = await tester.run_all_tests()
        
        print("\n" + "=" * 60)
        print("🎯 ТЕСТИРОВАНИЕ ЗАВЕРШЕНО")
        print("=" * 60)
        
        # Сохранение результатов в файл
        with open("test_results.json", "w", encoding="utf-8") as f:
            json.dump(results, f, indent=2, ensure_ascii=False)
            
        print("📄 Результаты сохранены в test_results.json")
        
    except Exception as e:
        print(f"💥 Критическая ошибка при запуске тестов: {e}")
        return 1
    
    return 0


if __name__ == "__main__":
    exit_code = asyncio.run(main())
    exit(exit_code)