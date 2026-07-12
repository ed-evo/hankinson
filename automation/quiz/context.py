import logging
import inspect

DRY_RUN: bool = False
LOG_LEVEL: str = "INFO"

def getLogger(logger_name: str | None = None):
    if not logger_name:
        frame = inspect.currentframe().f_back
        module_name = frame.f_globals.get('__name__')
        func_name = frame.f_code.co_name
        logger_name = f"{module_name}.{func_name}"
    logger = logging.getLogger(logger_name)

    numeric_level = getattr(logging, LOG_LEVEL.upper(), logging.INFO)
    
    logger.setLevel(numeric_level)

    console_handler = logging.StreamHandler()
    console_handler.setLevel(numeric_level)

    formatter = logging.Formatter(
        fmt="%(asctime)s | %(name)-12s | %(levelname)-8s | "
            f"{'[DRY-RUN] ' if DRY_RUN else ''}%(message)s",
        datefmt="%Y-%m-%d %H:%M:%S"
    )
    
    console_handler.setFormatter(formatter)
    logger.addHandler(console_handler)
    
    return logger