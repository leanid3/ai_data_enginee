"""Tool for generating Airflow DAG definitions."""

from __future__ import annotations

from datetime import datetime
from textwrap import dedent
from typing import Any, Dict, List, Optional

from croniter import croniter

from ..base import BaseTool, ToolInput, ToolOutput
from ..utils.validators import ensure_non_empty, ensure_unique


class AirflowEndpointConfig(ToolInput):
    type: str
    config: Dict[str, Any]


class TransformationConfig(ToolInput):
    name: str
    parameters: Dict[str, Any]


class AirflowNotifications(ToolInput):
    on_failure: Optional[str] = None
    on_success: Optional[str] = None


class AirflowDagInput(ToolInput):
    dag_id: str
    description: str
    schedule: str = "@daily"
    source: AirflowEndpointConfig
    target: AirflowEndpointConfig
    transformations: List[TransformationConfig] = []
    notifications: AirflowNotifications = AirflowNotifications()


class AirflowDagOutput(ToolOutput):
    dag_code: str
    dag_config: Dict[str, Any]
    operators: List[Dict[str, Any]]
    requirements: List[str]
    deployment_instructions: str
    monitoring_dashboard: Dict[str, Any]


class CreateAirflowDagTool(BaseTool[AirflowDagInput, AirflowDagOutput]):
    name = "create_airflow_dag"
    description = "Generate an Airflow DAG for ETL processes"
    input_schema = AirflowDagInput
    output_schema = AirflowDagOutput

    async def execute(self, params: AirflowDagInput) -> AirflowDagOutput:
        ensure_non_empty(params.dag_id, "dag_id")
        croniter.expand(params.schedule)  # Validate schedule

        dag_code = self._render_dag_code(params)

        dag_config = {
            "dag_id": params.dag_id,
            "schedule_interval": params.schedule,
            "start_date": datetime.utcnow().isoformat(),
            "catchup": False,
            "max_active_runs": 1,
            "description": params.description,
        }

        operators = self._build_operator_metadata(params)
        requirements = ["apache-airflow", "pendulum"]
        deployment_instructions = "Place DAG file under airflow/dags and restart the scheduler"
        monitoring_dashboard = {
            "type": "default",
            "url": f"/admin/airflow/graph?dag_id={params.dag_id}",
        }

        return AirflowDagOutput(
            success=True,
            dag_code=dag_code,
            dag_config=dag_config,
            operators=operators,
            requirements=requirements,
            deployment_instructions=deployment_instructions,
            monitoring_dashboard=monitoring_dashboard,
        )

    def _render_dag_code(self, params: AirflowDagInput) -> str:
        transformations_code = "\n".join(
            self._render_transformation(task, index)
            for index, task in enumerate(params.transformations, start=1)
        )

        dag_template = dedent(
            f"""
            from airflow import DAG
            from airflow.operators.python import PythonOperator
            from datetime import datetime


            def extract(**context):
                # Implement extraction from {params.source.type}
                return context


            def load(**context):
                # Implement load to {params.target.type}
                return context


            default_args = {{
                "owner": "data-eng",
                "depends_on_past": False,
                "retries": 1,
            }}


            with DAG(
                dag_id="{params.dag_id}",
                default_args=default_args,
                description="{params.description}",
                schedule_interval="{params.schedule}",
                start_date=datetime.utcnow(),
                catchup=False,
                tags=["generated", "etl"],
            ) as dag:
                extract_task = PythonOperator(
                    task_id="extract",
                    python_callable=extract,
                )

                {transformations_code}

                load_task = PythonOperator(
                    task_id="load",
                    python_callable=load,
                )

                extract_task >> [{', '.join(f'transform_{i}' for i in range(1, len(params.transformations) + 1))}] >> load_task
            """
        ).strip()

        return dag_template

    def _render_transformation(self, task: TransformationConfig, index: int) -> str:
        return dedent(
            f"""
                transform_{index} = PythonOperator(
                    task_id="transform_{index}",
                    python_callable=lambda **context: context,
                )
            """
        )

    def _build_operator_metadata(self, params: AirflowDagInput) -> List[Dict[str, Any]]:
        operators = [
            {"task_id": "extract", "operator_type": "PythonOperator", "dependencies": []}
        ]
        for index, _ in enumerate(params.transformations, start=1):
            operators.append(
                {
                    "task_id": f"transform_{index}",
                    "operator_type": "PythonOperator",
                    "dependencies": ["extract"],
                }
            )
        dependencies = [f"transform_{index}" for index in range(1, len(params.transformations) + 1)]
        operators.append(
            {
                "task_id": "load",
                "operator_type": "PythonOperator",
                "dependencies": dependencies,
            }
        )
        return operators


