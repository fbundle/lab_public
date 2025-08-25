from __future__ import annotations
from typing import Tuple, Iterable, Callable, Iterator

import ppft as pp


class Pool(pp.Server):
    def __enter__(self) -> Pool:
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.print_stats()
        self.destroy()

    def imap(self, func: Callable, args_iter: Iterable[Tuple], **kwargs) -> Iterator:
        task_list = []
        for args in args_iter:
            task = self.submit(func=func, args=args, **kwargs)
            task_list.append(task)

        for task in task_list:
            yield task.__call__()