from dataclasses import dataclass


@dataclass(frozen=True)
class CfgExample:
    ftp = """
[server]
host = sftp.example.com
port = 22
username = myuser
known_hosts = ~/.ssh/known_hosts  ; раскомментируй, если нужен

[auth]
password = mysecretpassword
; private_key = ~/.ssh/id_rsa      ; или используй ключ (раскомментируй один из двух)

[paths]
remote_dir = /upload/
    """
    llm = """"""

    def get_making_config(self, name):
        if name == "ftp":
            return self.ftp
        elif name == "llm":
            return self.llm
        else:
            raise ValueError(f"Unknown config name: {name}")


def make_cfg(name):
    return CfgExample().get_making_config(name)