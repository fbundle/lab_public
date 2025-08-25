"""
Code taken from https://github.com/twitter-research/graph-neural-pde/blob/main/src/data.py
"""
import os

import numpy as np
import torch
from torch_geometric import transforms as T
from torch_geometric.data import InMemoryDataset, Data
from torch_geometric.datasets import Planetoid, Amazon, Coauthor, WebKB, WikipediaNetwork
from torch_geometric.utils import to_undirected


def get_dataset(name: str = "Cora", use_lcc: bool = False) -> InMemoryDataset:
    path = os.path.abspath(f"data/{name}")
    if name in ['Cora', 'Citeseer', 'Pubmed']:
        dataset = Planetoid(path, name)
    elif name in ['Computers', 'Photo']:
        dataset = Amazon(path, name)
    elif name == 'CoauthorCS':
        dataset = Coauthor(path, 'CS')
    elif name in ['cornell', 'texas', 'wisconsin']:
        dataset = WebKB(root=path, name=name, transform=T.NormalizeFeatures())
    elif name in ['chameleon', 'squirrel']:
        dataset = WikipediaNetwork(root=path, name=name, transform=T.NormalizeFeatures())
    else:
        raise Exception('Unknown dataset.')

    if use_lcc:
        lcc = get_largest_connected_component(dataset)

        x_new = dataset.data.x[lcc]
        y_new = dataset.data.y[lcc]

        row, col = dataset.data.edge_index.numpy()
        edges = [[i, j] for i, j in zip(row, col) if i in lcc and j in lcc]
        edges = remap_edges(edges, get_node_mapper(lcc))

        data = Data(
            x=x_new,
            edge_index=torch.LongTensor(edges),
            y=y_new,
            train_mask=torch.zeros(y_new.size()[0], dtype=torch.bool),
            test_mask=torch.zeros(y_new.size()[0], dtype=torch.bool),
            val_mask=torch.zeros(y_new.size()[0], dtype=torch.bool)
        )
        dataset.data = data
    train_mask_exists = True
    try:
        dataset.data.train_mask
    except AttributeError:
        train_mask_exists = False

    if name == 'ogbn-arxiv':
        split_idx = dataset.get_idx_split()
        ei = to_undirected(dataset.data.edge_index)
        data = Data(
            x=dataset.data.x,
            edge_index=ei,
            y=dataset.data.y,
            train_mask=split_idx['train'],
            test_mask=split_idx['test'],
            val_mask=split_idx['valid'])
        dataset.data = data
        train_mask_exists = True

    return dataset


def get_node_mapper(lcc: np.ndarray) -> dict:
    mapper = {}
    counter = 0
    for node in lcc:
        mapper[node] = counter
        counter += 1
    return mapper


def get_largest_connected_component(dataset: InMemoryDataset) -> np.ndarray:
    remaining_nodes = set(range(dataset.data.x.shape[0]))
    comps = []
    while remaining_nodes:
        start = min(remaining_nodes)
        comp = get_component(dataset, start)
        comps.append(comp)
        remaining_nodes = remaining_nodes.difference(comp)
    return np.array(list(comps[np.argmax(list(map(len, comps)))]))


def remap_edges(edges: list, mapper: dict) -> list:
    row = [e[0] for e in edges]
    col = [e[1] for e in edges]
    row = list(map(lambda x: mapper[x], row))
    col = list(map(lambda x: mapper[x], col))
    return [row, col]


def set_train_val_test_split(
        seed: int,
        data: Data,
        num_development: int = 1500,
        num_per_class: int = 20) -> Data:
    rnd_state = np.random.RandomState(seed)
    num_nodes = data.y.shape[0]
    development_idx = rnd_state.choice(num_nodes, num_development, replace=False)
    test_idx = [i for i in np.arange(num_nodes) if i not in development_idx]

    train_idx = []
    rnd_state = np.random.RandomState(seed)
    for c in range(data.y.max() + 1):
        class_idx = development_idx[np.where(data.y[development_idx].cpu() == c)[0]]
        train_idx.extend(rnd_state.choice(class_idx, num_per_class, replace=False))

    val_idx = [i for i in development_idx if i not in train_idx]

    def get_mask(idx):
        mask = torch.zeros(num_nodes, dtype=torch.bool)
        mask[idx] = 1
        return mask

    data.train_mask = get_mask(train_idx)
    data.val_mask = get_mask(val_idx)
    data.test_mask = get_mask(test_idx)

    return data


def get_component(dataset: InMemoryDataset, start: int = 0) -> set:
    visited_nodes = set()
    queued_nodes = set([start])
    row, col = dataset.data.edge_index.numpy()
    while queued_nodes:
        current_node = queued_nodes.pop()
        visited_nodes.update([current_node])
        neighbors = col[np.where(row == current_node)[0]]
        neighbors = [n for n in neighbors if n not in visited_nodes and n not in queued_nodes]
        queued_nodes.update(neighbors)
    return visited_nodes
