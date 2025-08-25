from __future__ import annotations

import math
import random

import torch
from torch import nn
from torch_scatter import scatter_sum


# pip install pyg_lib torch_scatter torch_sparse torch_cluster torch_spline_conv -f https://data.pyg.org/whl/torch-2.0.0+cpu.html

def fully_connected_graph(n: int) -> torch.Tensor:
    e = torch.tensor([
        n * list(range(n)),
        sum([[i] * n for i in range(n)], [])
    ], dtype=torch.int64)
    return e


def sparse_multihead_softmax(values: torch.Tensor, index: torch.Tensor, eps: float = 1e-16) -> torch.Tensor:
    """
    :param values: (b, m, h)
    :param index: (m,)
    :param eps:
    :return: (b, m, h)
    """
    assert len(values.shape) == 3
    assert len(index.shape) == 1
    assert index.dtype == torch.int64
    assert values.shape[1] == index.shape[0]

    values = values - values.max()  # shift does not change softmax, use this to avoid overflow and reduce numerical err
    exp_logits = torch.exp(values)
    exp_logits_sum = scatter_sum(src=exp_logits, index=index, dim=1)[:, index, :]
    return exp_logits / (exp_logits_sum + eps)


def sparse_multihead_attention(
        q: torch.Tensor, k: torch.Tensor, e: torch.Tensor, h: torch.Tensor,
):
    """
    :param q: (b, n, d) - (batch, node, dim)
    :param k: (b, n, d) - (batch, node, dim)
    :param e: (2, m) - (2, edge)
    :param h: (d,) - (dim,)
    :return: (b, m, h) - (batch, edge, head)
    """
    assert len(q.shape) == 3
    assert len(k.shape) == 3
    assert q.shape == k.shape

    assert len(e.shape) == 2
    assert e.shape[0] == 2
    assert e.dtype == torch.int64

    assert len(h.shape) == 1
    assert h.shape[0] == q.shape[2]

    src_q = q[:, e[0, :], :]  # (b, m, d)
    dst_k = k[:, e[1, :], :]  # (b, m, d)
    a_e = src_q * dst_k  # (b, m, d)
    a_e_sum = scatter_sum(src=a_e, index=h, dim=2)  # (b, m, h)
    return a_e_sum


class SparseMultiheadAttention(nn.Module):
    def __init__(self,
                 key_dim: int,
                 hidden_dim: int,
                 num_heads: int = 1,
                 ):
        assert hidden_dim % num_heads == 0

        super().__init__()

        self.num_heads = num_heads
        self.key_dim = key_dim
        self.hidden_dim = hidden_dim

        self.register_buffer("h", torch.tensor(
            [i // (hidden_dim // num_heads) for i in range(hidden_dim)]
            , dtype=torch.int64,
        ))

    def __repr__(self) -> str:
        return f"{self.__class__.__name__}(key_dim={self.key_dim}, hidden_dim={self.hidden_dim}, num_heads={self.num_heads})"

    def forward(self, q: torch.Tensor, k: torch.Tensor, e: torch.Tensor | None = None,
                r: torch.Tensor | None = None) -> torch.Tensor:
        """
        :param q: (b, n, d)
        :param k: (b, n, d)
        :param e: (2, m)
        :param r: (m, )
        :return: (b, m, h)
        """
        assert len(q.shape) == 3
        assert len(k.shape) == 3
        assert q.shape == k.shape
        if e is None:  # fully connected graph - normal multihead attention
            n = q.shape[1]
            e = fully_connected_graph(n).to(self.device)

        assert len(e.shape) == 2
        assert e.shape[0] == 2
        assert e.dtype == torch.int64
        if r is None:
            r = e[0, :]
        assert len(r.shape) == 1
        assert r.shape[0] == e.shape[1]

        a_values = sparse_multihead_attention(
            q=q,
            k=k,
            e=e,
            h=self.h,
        ) / math.sqrt(self.key_dim)  # (b, m, h)
        a_values = sparse_multihead_softmax(values=a_values, index=r)
        return a_values


if __name__ == "__main__":
    def _attention(q: torch.Tensor, k: torch.Tensor):
        """
        :param q: (b, n, d)
        :param k: (b, n, d)
        :return: (b, n, n)
        """
        assert len(q.shape) == 3
        assert len(k.shape) == 3
        assert q.shape == k.shape
        a_sum = torch.bmm(q, k.transpose(1, 2))
        return a_sum


    def _sparse_multihead_attention_elementwise(
            q: torch.Tensor, k: torch.Tensor, e: torch.Tensor, h: torch.Tensor,
    ):
        """
        :param q: (b, n, d) - (batch, node, dim)
        :param k: (b, n, d) - (batch, node, dim)
        :param e: (2, m) - (2, edge)
        :param h: (d,) - (dim,)
        :return: (b, m, h) - (batch, edge, head)
        """
        a_sum_list = []
        for i in h.unique():
            q_i = q[:, :, h == i]
            k_i = k[:, :, h == i]
            a_sum_i = _attention(q_i, k_i)
            a_sum_list.append(a_sum_i)
        a_sum = torch.stack(a_sum_list, dim=3)  # (b, n, n, h)
        a_e_sum_list = []
        for k in range(e.shape[1]):
            i, j = e[:, k]
            a_e_sum_list.append(a_sum[:, i, j, :])
        a_e_sum = torch.stack(a_e_sum_list, dim=1)  # (b, m, h)
        return a_e_sum


    # test sparse_multi_headed_attention
    b, n, d = 3, 5, 4
    q = torch.rand(b, n, d)
    k = torch.rand(b, n, d)
    e = []
    for i in range(n):
        for j in range(n):
            if random.random() < 0.5:
                e.append((i, j))
    e = torch.tensor(e).T
    h = torch.tensor([0, 0, 1, 1])
    expected = _sparse_multihead_attention_elementwise(q, k, e, h)
    actual = sparse_multihead_attention(q, k, e, h)
    assert torch.isclose(expected, actual).all()
