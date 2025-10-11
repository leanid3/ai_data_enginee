from dataclasses import dataclass
from pathlib import Path
from typing import Optional


@dataclass(frozen=True)
class SFTPConfig:
    host: str
    port: int
    username: str
    password: Optional[str]
    private_key: Optional[Path]
    remote_dir: str
    known_hosts: Optional[str]  
    
    def __post_init__(self):
        if not self.password and not self.private_key:
            raise ValueError("Требуется password или private_key")
        if self.private_key:
            object.__setattr__(self, 'private_key', Path(self.private_key))
            