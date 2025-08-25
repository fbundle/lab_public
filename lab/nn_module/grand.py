from typing import Tuple

import torch
import torch_scatter
import torch_sparse
from torch import Tensor, nn

from .attention import SparseMultiheadAttention, fully_connected_graph
from .neural_ode import ODEFunc


def sparse_batch_matmul(values: torch.Tensor, edge: torch.Tensor, dim: Tuple[int, int],
                        dense: torch.Tensor) -> torch.Tensor:
    """
    :param values: (*, m)
    :param edge: (2, m)
    :param dim: (n1, n2)
    :param dense: (*, n2, n3)
    :return: sparse @ dense (*, n1, n3)
    """
    assert len(values.shape) >= 2
    assert len(edge.shape) == 2
    assert edge.shape[0] == 2
    assert values.shape[-1] == edge.shape[1]
    assert len(dense.shape) >= 3
    assert dense.shape[:-2] == values.shape[:-1]
    assert dim[1] == dense.shape[-2]

    return torch_sparse.spmm(index=edge, value=values, m=dim[0], n=dim[1], matrix=dense)


class LinearGRAND(ODEFunc):
    def __init__(self,
                 num_nodes: int,
                 diffusion_dim: int,
                 attention_dim: int,
                 edge_index: torch.Tensor | None = None,
                 num_heads: int = 1,
                 softmax_dim: int = 1,
                 bias: bool = False,
                 ):
        """
        :param edge_index: (2, m)
        """
        super().__init__()

        if edge_index is None:  # fully connected graph - normal multihead attention
            edge_index = fully_connected_graph(num_nodes)

        self.num_nodes = num_nodes
        self.register_buffer("edge_index", edge_index)
        self.register_parameter("w_q", nn.Parameter(torch.empty(diffusion_dim, attention_dim, dtype=torch.float32)))
        self.register_parameter("w_k", nn.Parameter(torch.empty(diffusion_dim, attention_dim, dtype=torch.float32)))

        self.add_module("attention", SparseMultiheadAttention(
            key_dim=diffusion_dim,
            hidden_dim=attention_dim,
            num_heads=num_heads,
        ))

        self.attention_values = None
        self.reduce_operator = None
        self.radius2 = None

        self.softmax_dim = softmax_dim

        self._reset_parameters()

    def _reset_parameters(self):
        nn.init.xavier_uniform_(self.w_q)
        nn.init.xavier_uniform_(self.w_k)

    def on_before_integration(self, x0: torch.Tensor, t_list: torch.Tensor):
        """
        call before every integration
        """
        self.attention_values = self.attention.forward(
            q=x0 @ self.w_q,
            k=x0 @ self.w_k,
            e=self.edge_index,
            r=self.edge_index[1 - self.softmax_dim, :],
        ).mean(dim=-1)  # average over all heads

        b, n, _ = x0.shape
        self.reduce_operator = torch_scatter.scatter_sum(
            src=self.attention_values,
            index=self.edge_index[1 - self.softmax_dim, :],
        ).view(b, n, 1)

    def forward(self, t: float, x: Tensor):
        """
        :param t:
        :param x: (b, n, d)
        """
        assert len(x.shape) == 3

        f = sparse_batch_matmul(
            self.attention_values,
            self.edge_index,
            (self.num_nodes, self.num_nodes),
            x,
        ) - self.reduce_operator * x

        return f
