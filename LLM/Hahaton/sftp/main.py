# from config.FTP import load_ftp_cfg
# from sftp.device import SFTPDevice
# import asyncio
# from pathlib import Path
#
#
# async def main():
#     config = load_ftp_cfg("../config/FTP/ftp_config.cfg")
#     device = SFTPDevice(config)
#
#     async with device.connect() as client:
#         files = await client.listdir(config.remote_dir)
#         print("Files:", files)
#
# if __name__ == "__main__":
#     asyncio.run(main())