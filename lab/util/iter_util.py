from __future__ import annotations

import functools
from typing import Iterable, Callable, Any, Iterator


class Iter:
    def __init__(self, it: Iterable):
        self.it: Iterable | Iterator = it

    def __iter__(self):
        return Iter(it=iter(self.it))

    def __next__(self) -> Any:
        return next(self.it)

    def map(self, f: Callable[[Any], Any]) -> Iter:
        return Iter(it=map(f, self.it))

    def filter(self, f: Callable[[Any], bool]) -> Iter:
        return Iter(it=filter(f, self.it))

    def reduce(self, f: Callable[[Any, Any], Any], initial: Any = object()) -> Any:
        functools.reduce(function=f, sequence=self.it, initial=initial)

    def flat_map(self, f: Callable[[Any], Iterable[Any]]) -> Iter:
        def helper() -> Iterator:
            for i in self.it:
                yield from f(i)

        return Iter(it=helper())
