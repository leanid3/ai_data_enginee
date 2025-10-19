from enum import Enum
from typing import List


class AgentRole(str, Enum):
    ROUTER = "router"
    DATA_ANALYZER = "data_analyzer"
    DB_SELECTOR = "db_selector"
    DDL_GENERATOR = "ddl_generator"
    ETL_BUILDER = "etl_builder"
    QUERY_OPTIMIZER = "query_optimizer"
    REPORT_GENERATOR = "report_generator"
    CONSOLIDATOR = "consolidator"


class AgentChain:
    CHAINS = {
        "full_process": [
            AgentRole.DATA_ANALYZER,
            AgentRole.DB_SELECTOR,
            AgentRole.DDL_GENERATOR,
            AgentRole.ETL_BUILDER,
            AgentRole.QUERY_OPTIMIZER,
            AgentRole.REPORT_GENERATOR,
        ],
        "analyze": [AgentRole.DATA_ANALYZER, AgentRole.REPORT_GENERATOR],
        "recommend": [
            AgentRole.DATA_ANALYZER,
            AgentRole.DB_SELECTOR,
            AgentRole.REPORT_GENERATOR,
        ],
        "generate_ddl": [
            AgentRole.DATA_ANALYZER,
            AgentRole.DB_SELECTOR,
            AgentRole.DDL_GENERATOR,
        ],
        "build_pipeline": [
            AgentRole.ETL_BUILDER,
            AgentRole.QUERY_OPTIMIZER,
            AgentRole.REPORT_GENERATOR,
        ],
        "optimize": [AgentRole.QUERY_OPTIMIZER, AgentRole.REPORT_GENERATOR],
    }

    @classmethod
    def get_chain(cls, operation_type: str) -> List[AgentRole]:
        return cls.CHAINS.get(operation_type, cls.CHAINS["analyze"])