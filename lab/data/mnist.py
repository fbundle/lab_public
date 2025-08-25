from typing import *

import torch
import torchvision


def get_data() -> Tuple[torch.Tensor, torch.Tensor]:
    dataset = torchvision.datasets.MNIST("db", download=True)
    X, y = dataset.data, dataset.targets
    return X, y
