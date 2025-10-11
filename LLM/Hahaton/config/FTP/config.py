import configparser
from .struct.sftp_conf_struct import SFTPConfig
from pathlib import Path

class ConfigWorker:
    def __init__(self, path: str):
        self.path = path
        self.cfg = configparser.ConfigParser()
    
    def __call__(self) -> SFTPConfig:
        self.cfg.read(self.path)
        known_hosts = self.cfg["server"].get("known_hosts")
        if known_hosts:
            known_hosts = str(Path(known_hosts).expanduser().resolve())
        return SFTPConfig(
            host=self.cfg["server"]["host"],
            port=int(self.cfg["server"]["port"]),
            username=self.cfg["server"]["username"],
            password=self.cfg["auth"].get("password"),
            private_key=self.cfg["auth"].get("private_key"),
            remote_dir=self.cfg["paths"]["remote_dir"],
            known_hosts=self.cfg["server"]["known_hosts"]
        )
        

def load_ftp_cfg(path):
    return ConfigWorker(path)()
        
"""
Use:
from config import load_ftp_cfg

config = load_ftp_cfg(path_to_config)
host = config.host
...
"""
