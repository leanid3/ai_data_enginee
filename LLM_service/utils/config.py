from dataclasses import dataclass
from pathlib import Path
import os

from cfg_examples import make_cfg

CFG_PATH = "../config"


@dataclass
class ConfigPaths:
    ftp: Path = Path(CFG_PATH, "FTP")
    llm: Path = Path(CFG_PATH, "LLM")
    api: Path = Path(CFG_PATH, "API")

    def __post_init__(self):
        self.create_config_files()

    def create_config_files(self):
        """Создает все нужные конфигурационные файлы"""
        self.create_config(self.ftp, "ftp_config.cfg", get_config("ftp"))
        llm_content = get_config("llm")
        if llm_content.strip():
            self.create_config(self.llm, "llm_config.cfg", llm_content)

    @staticmethod
    def create_config(path: Path, filename: str, content: str):
        """Создает конфигурационный файл с указанным содержимым"""
        try:
            path.mkdir(parents=True, exist_ok=True)
            file_path = path / filename
            with open(file_path, "w", encoding="utf-8") as f:
                f.write(content.strip())  # убираем лишние пробелы
            print(f"Создан конфигурационный файл: {file_path}")
        except Exception as e:
            raise Exception(f"Cannot create config file: {e}")


class GetConfig:
    def __init__(self):
        self.config_paths = ConfigPaths()

    def get_config_path(self, config_type: str) -> Path:
        """Возвращает путь к директории с конфигурацией"""
        config_map = {
            "ftp": self.config_paths.ftp,
            "llm": self.config_paths.llm,
            "api": self.config_paths.api
        }
        if config_type not in config_map:
            raise ValueError(f"Unknown config type: {config_type}")

        return config_map[config_type]

    def get_config_file_path(self, config_type: str) -> Path:
        """Возвращает полный путь к конфигурационному файлу"""
        base_path = self.get_config_path(config_type)
        return base_path / f"{config_type}_config.ini"


def get_config():
    return GetConfig()

"""
Use:
# Создаем все конфиги автоматически
config_manager = get_cfg()

# Получаем пути к конфигурационным файлам
ftp_config_path = config_manager.get_config_file_path("ftp")
print(f"FTP config location: {ftp_config_path}")

# Проверяем, что файл создан
if ftp_config_path.exists():
    print("FTP config file created successfully!")
    
    # Читаем содержимое
    with open(ftp_config_path, 'r') as f:
        content = f.read()
        print("File content:")
        print(content)
"""