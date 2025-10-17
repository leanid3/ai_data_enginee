import logging
import inspect
import queue
from logging import Logger
from logging.handlers import RotatingFileHandler, QueueHandler, QueueListener
from pathlib import Path


class ColoredFormatter(logging.Formatter):
    COLORS = {
        'DEBUG': '\033[2;36m',
        'INFO': '\033[0;32m',
        'WARNING': '\033[0;33m',
        'ERROR': '\033[0;31m',
        'CRITICAL': '\033[1;31m'
    }
    RESET = '\033[0m'

    def format(self, record):
        log_color = self.COLORS.get(record.levelname, self.RESET)
        record.levelname = f"{log_color}{record.levelname}{self.RESET}"
        return super().format(record)


class AutoClassFuncLogger(Logger):
    """
    Cls and func auto-detecting logger
    Add extra={'cls': ..., 'func': ...} for all msg`s.
    """
    def _log(self, level, msg, args, exc_info=None, extra=None, stack_info=False, stacklevel=1):
        if extra is None:
            extra = {}

        frame = inspect.currentframe()
        try:
            for _ in range(3):
                if frame.f_back:
                    frame = frame.f_back
                else:
                    break

            func_name = frame.f_code.co_name
            cls_name = "Unknown"

            if 'self' in frame.f_locals:
                cls_name = frame.f_locals['self'].__class__.__name__
            elif 'cls' in frame.f_locals:
                cls_name = frame.f_locals['cls'].__name__

            extra.update(cls=cls_name, func=func_name)
        finally:
            del frame 

        super()._log(level, msg, args, exc_info, extra, stack_info, stacklevel)


logging.setLoggerClass(AutoClassFuncLogger)


def setup_logger(
    version: str,
    name: str = 'DataEngineer',
    log_file: str = None,
    level: int = logging.INFO
) -> AutoClassFuncLogger:
    logger = logging.getLogger(name)
    if logger.handlers:
        return logger

    base_format = f'{version} - %(cls)s - %(func)s - %(asctime)s - %(name)s - %(levelname)s - %(message)s'

    file_formatter = logging.Formatter(base_format)
    console_formatter = ColoredFormatter(base_format)

    console_handler = logging.StreamHandler()
    console_handler.setFormatter(console_formatter)

    file_handler = None
    if log_file:
        log_path = Path(log_file)
        log_path.parent.mkdir(parents=True, exist_ok=True)
        file_handler = RotatingFileHandler(
            log_file,
            maxBytes=10 * 1024 * 1024,
            backupCount=5,
            encoding='utf-8'
        )
        file_handler.setFormatter(file_formatter)

    log_queue = queue.Queue(-1)
    queue_handler = QueueHandler(log_queue)
    queue_handler.setLevel(level)

    handlers = [console_handler]
    if file_handler:
        handlers.append(file_handler)

    listener = QueueListener(log_queue, *handlers, respect_handler_level=True)
    listener.start()

    logger.listener = listener
    logger.addHandler(queue_handler)
    logger.setLevel(level)
    logger.propagate = False

    return logger


def get_logger(version: str) -> AutoClassFuncLogger:
    """Getting global logger"""
    log_file = 'logs/system.log'
    return setup_logger(version, 'DataEngineer', log_file)


"""
Use:
logger = get_logger(version)
logger.info(...) # Using queue
"""