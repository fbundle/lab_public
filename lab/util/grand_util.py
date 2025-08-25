from typing import Tuple

import numpy as np


def smoothness(x: np.ndarray) -> Tuple[np.ndarray, np.ndarray]:
    """
    :param x: (b, n, d)
    project all columns of (n, d) into 1 vector.
    smoothness is the ratio between norm2 of the projection and the norm2 of the original vector

    """
    norm_2 = (x ** 2).sum(axis=1)  # (b, d)
    f0_norm_2 = x.sum(axis=1)**2 / x.shape[1]
    return f0_norm_2 / norm_2, norm_2


def normalize(x: np.ndarray) -> np.ndarray:
    x = x - x.mean(axis=1, keepdims=True)  # mean 0
    x = x / x.std(axis=1, keepdims=True)  # std 1
    return x
