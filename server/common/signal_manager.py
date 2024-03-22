import signal

class SignalManagerMeta(type):
    _instances = {}

    def __call__(cls, *args, **kwargs):
        if cls not in cls._instances:
            cls._instances[cls] = super().__call__(*args, **kwargs)
        return cls._instances[cls]

class SignalManager(metaclass=SignalManagerMeta):
    SIGTERM = signal.SIGTERM
    _instance = None

    def __init__(self):
        self._handlers_by_signum = {}

    def _reduce_handlers(handlers):
        return lambda signum, frame: [handler(signum, frame) for handler in handlers]

    def add_handler(self, signum, handler):
        if (self._handlers_by_signum.get(signum) is None):
            self._handlers_by_signum[signum] = [handler]
            signal.signal(signum, SignalManager._reduce_handlers(self._handlers_by_signum[signum]))
        else:
            self._handlers_by_signum[signum].append(handler)
