from __future__ import annotations


class NamedTuple:
    def __init__(self, **kwargs):
        for field_name, field_type_str in self.__annotations__.items():
            setattr(self, field_name, kwargs.get(field_name, None))


class Function(NamedTuple):
    def evaluate(self, *args: float) -> float:
        raise NotImplementedError

    def partial_differentiate(self, i: int) -> Function:
        raise NotImplementedError


class Variable(Function):
    name: str

    def evaluate(self, *args: float) -> float:
        return args[0]

    def partial_differentiate(self, i: int) -> Function:
        assert i == 1
        return Constant(value=1)


class Constant(Function):
    value: float

    def evaluate(self, *args: float) -> float:
        return self.value

    def partial_differentiate(self, i: int) -> Function:
        return Constant(value=0)
