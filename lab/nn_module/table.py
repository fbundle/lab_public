import torch
from torch import nn


class Table(nn.Module):
    """
    Table: table of size (d1, d2)
    """

    def __init__(self, d1: int, d2: int):
        super().__init__()
        self.register_parameter("table", nn.Parameter(torch.empty(d1, d2)))
        self._reset_parameters()

    def _reset_parameters(self):
        nn.init.xavier_uniform_(self.table)

    def __repr__(self):
        return f"{self.__class__.__name__}(d1={self.table.shape[0]}, d2={self.table.shape[1]})"

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        """
        :param x: (*, d1)
        :return: (*, d1, d2)
        """
        assert x.shape[-1] == self.table.shape[0]
        x_ = x.unsqueeze(2)  # (*, d1) -> (*, d1, 1)
        t_ = self.table.unsqueeze(0)  # (d1, d2) -> (1, d1, d2)
        s = x_ * t_ # (*, d1, d2)
        return s


if __name__ == "__main__":
    t = Table(5, 3)
    print(t)
    print(t.table)
    x = torch.Tensor([
        [0, 1, 1, 0, 0],
        [1, 0, 1, 1, 0],
    ])
    print(t.forward(x))
