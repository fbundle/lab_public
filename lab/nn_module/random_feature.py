import torch
from torch import nn


class RandomFeature(nn.Module):
    def __init__(self, dim: int, width: int):
        super().__init__()
        r = torch.randn(dim, width)
        r = r / torch.linalg.norm(r, dim=0)
        self.register_buffer("rotation", r)

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        """
        :param x: (n, dim)
        :return: (n, width)
        """
        return x @ self.rotation
