from typing import *

import torch
from torch import nn


class Functional(nn.Module):
    """
    Functional: wrapper for function
    """

    def __init__(self, f: Callable, name: Optional[str] = None):
        super().__init__()
        if name is None:
            name = f"{f}"
        self.f = f
        self.name = name

    def __repr__(self):
        return f"{self.__class__.__name__} {self.name}"

    def forward(self, *args, **kwargs):
        return self.f(*args, **kwargs)


def make_mlp(dim_list: List[int], act: Callable[[], nn.Module] = nn.Tanh) -> List[nn.Module]:
    assert len(dim_list) >= 2
    sequence: List[nn.Module] = []

    for i, (dim_in, dim_out) in enumerate(zip(dim_list[:-1], dim_list[1:])):
        sequence.append(nn.Linear(in_features=dim_in, out_features=dim_out))
        if i < len(dim_list) - 2:  # not last layer
            sequence.append(act())

    return sequence


class Residual(nn.Module):
    def __init__(self, *module_list: nn.Module):
        super().__init__()
        self.module_list = module_list
        for i, module in enumerate(self.module_list):
            self.register_module(f"{i}", module)

    def forward(self, x0: torch.Tensor):
        x = x0
        for module in self.module_list:
            x = x + module.forward(x)
        return x


