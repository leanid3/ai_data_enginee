import asyncio
import asyncssh
from typing import Optional, AsyncIterator
from contextlib import asynccontextmanager
import logging

from config.FTP import SFTPConfig
from sftp.client.cleint import SFTPClientAsync


class SFTPDevice:
    def __init__(self, logger:logging.Logger, config: SFTPConfig, max_retr: int = 3, retr_delay: float = 1.0):
        self.config = config
        self.max_retr = max_retr
        self.retr_delay = retr_delay
        self.logger = logger

    @asynccontextmanager
    async def connect(self) -> AsyncIterator[SFTPClientAsync]:
        conn: Optional[asyncssh.SSHClientConnection] = None
        sftp: Optional[asyncssh.SFTPClient] = None

        self.logger.info("Trying to connect to SFTP server: loading params")

        host = self.config.host
        port = self.config.port
        username = self.config.username
        password = self.config.password
        private_key = self.config.private_key
        known_hosts = self.config.known_hosts

        last_exception = None
        for attempt in range(1, self.max_retr + 1):
            try:
                self.logger.info(f"Connecting to {username}@{host}:{port} (attempt {attempt})")
                
                kh = asyncssh.read_known_hosts(known_hosts) if known_hosts else None

                conn = await asyncssh.connect(
                    host=host,
                    port=port,
                    username=username,
                    password=password,
                    client_keys=[private_key] if private_key else None,
                    known_hosts=kh,
                )
                sftp = await conn.start_sftp_client()
                self.logger.info("Connected successfully")
                break
            except Exception as e:
                last_exception = e
                self.logger.warning(f"Attempt {attempt} failed: {e}")
                if attempt < self.max_retr:
                    await asyncio.sleep(self.retr_delay)
        else:
            self.logger.error("All connection attempts have been exhausted")
            raise last_exception

        client = SFTPClientAsync(sftp)

        try:
            yield client
        except Exception as e:
            self.logger.error(f"Error when working with sFTP: {e}")
            raise
        finally:
            if conn and not conn.closed:
                conn.close()
                await conn.wait_closed()
                self.logger.info("Connection closed")

"""
Use:
from config.config import load_ft
from sftp.device import SFTPDevice
import asyncio

async def main():
    config = load_ftp_cfg("config/ftp_config.cfg")
    device = SFTPDevice(config)

    async with device.connect() as client:
        files = await client.listdir(config.remote_dir)
        print("Files:", files)

if __name__ == "__main__":
    asyncio.run(main())
"""