import asyncssh
from typing import Optional


class SFTPClientAsync:
    def __init__(self, sftp: asyncssh.SFTPClient):
        self._sftp = sftp

    async def upload(self, local_path: str, remote_path: str) -> None:
        await self._sftp.put(local_path, remote_path)

    async def download(self, remote_path: str, local_path: str) -> None:
        await self._sftp.get(remote_path, local_path)

    async def listdir(self, path: str = ".") -> list:
        return await self._sftp.readdir(path)

    async def exists(self, path: str) -> bool:
        try:
            await self._sftp.stat(path)
            return True
        except asyncssh.SFTPError:
            return False

    async def remove(self, path: str) -> None:
        await self._sftp.remove(path)