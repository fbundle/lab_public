from __future__ import annotations

from typing import *

from pathos.multiprocessing import ProcessPool


class Executor:
    global_pool: ProcessPool = ProcessPool()

    def __init__(self, f: Callable, pool: Optional[ProcessPool] = None):
        if pool is None:
            pool = Executor.global_pool
        self.pool = pool
        self.callable = f

    def imap(self, *args, **kwargs) -> Iterator:
        return self.pool.imap(self.callable, *args, **kwargs)

    def __call__(self, *args, **kwargs) -> Any:
        return self.callable(*args, **kwargs)


if __name__ == "__main__":
    @Executor
    def add(x: int, y: int) -> int:
        return x + y


    assert 4 == add(1, 3)

    for out in add.imap([1, 2, 3, 4], [2, 3, 4, 5]):
        print(out)
