from typing import Tuple, Iterable, Optional, Any, List, Callable

import numpy as np
import scipy as sp
import scipy.sparse.linalg


def gram_schmidt(V: np.ndarray, v: np.ndarray, Q: Optional[np.ndarray] = None):
    """
    :param V: (n, m) orthonormal columns
    :param v: (n,) vector
    :param Q: (n, n) Q-orthogonal
    :return:
        q: orthogonal to columns of Q
        r: projection on columns of Q
    """
    assert len(V.shape) == 2 and len(v.shape) == 1
    assert V.shape[0] == v.shape[0]
    n, m = V.shape

    if Q is None:
        Q = np.identity(n)

    if m == 0:
        return v, np.zeros(shape=(m,))
    r = V.T @ (Q @ v)
    q = v - (r.reshape(1, m) * V).sum(axis=1)
    return q, r


def qr_decomposition(A: np.ndarray) -> Tuple[np.ndarray, np.ndarray]:
    """
    :param A: (n, n) matrix
    :return:
        Q: (n, n) orthonormal
        R: (n, n) upper triangular
        A = QR
    """
    assert len(A.shape) == 2 and A.shape[0] == A.shape[1]
    n = A.shape[0]
    Q = np.zeros(shape=(n, n))
    R = np.zeros(shape=(n, n))
    for i in range(n):
        a_i = A[:, i]
        # gram-schmidt
        q_i, r_i = gram_schmidt(Q, a_i)
        # normalize
        q_i_norm = np.linalg.norm(q_i)
        r_i[i], q_i = q_i_norm, q_i / q_i_norm
        # assignment
        Q[:, i] = q_i
        R[:, i] = r_i
    return Q, R


def qr_algorithm(A: np.ndarray) -> Iterable[np.ndarray]:
    """
    :param A:
    :return: sequence of (A_t, U_t).
        A_t: (n, n) approaches upper triangular and similar to A
        U_t: (n, n) orthonormal
        A = U_t A_t U_t^T
    """
    assert len(A.shape) == 2 and A.shape[0] == A.shape[1]
    n = A.shape[0]
    U_t = np.identity(n)
    A_t = A
    while True:
        Q, R = qr_decomposition(A_t)
        A_t = R @ Q  # = Q^T Q R Q = Q^T A Q
        U_t = Q.T @ U_t
        yield A_t, U_t


def lanczos_iteration(A: np.ndarray, m: Optional[int] = None, eps: float=1e-1, stable: bool = False):
    """
    :param A: (n, n) real symmetric
    :param m: m < n
    :param stable:
    :return:
        T: (m, m) tridiagonal form
        V (n, m) orthonormal columns
        T = V^T A V
    """
    assert len(A.shape) == 2 and A.shape[0] == A.shape[1]
    assert np.isclose(A, A.T).all()
    
    n = A.shape[0]
    if m is None:
        m = n
    V = np.zeros(shape=(n, n))
    T = np.zeros(shape=(n, n))
    # v_0
    v_0 = np.random.normal(size=(n,))
    v_0 /= np.linalg.norm(v_0)
    w_0_ = A @ v_0
    alpha_0 = np.dot(w_0_, v_0)
    w_0 = w_0_ - alpha_0 * v_0

    V[:, 0] = v_0
    T[0, 0] = alpha_0

    v_prev, w_prev = v_0, w_0
    for j in range(1, m):
        beta_j = np.linalg.norm(w_prev)
        if beta_j > eps:
            v_j = w_prev / beta_j
            if stable:
                # added gram-schmidt here for numerical stability
                # not neccessary theoretically
                v_j, _ = gram_schmidt(V, v_j)
                v_j /= np.linalg.norm(v_j)
        else:
            v_j = np.random.normal(size=(n,))
            v_j, _ = gram_schmidt(V, v_j)
            v_j /= np.linalg.norm(v_j)

        w_j_ = A @ v_j
        alpha_j = np.dot(w_j_, v_j)
        w_j = w_j_ - alpha_j * v_j - beta_j * v_prev

        V[:, j] = v_j
        T[j, j] = alpha_j
        T[j, j - 1] = beta_j
        T[j - 1, j] = beta_j
        v_prev, w_prev = v_j, w_j
    return T[:m, :m], V[:, :m]


def lanczos_qr_eigh(A: np.ndarray, eps: float = 1e-3) -> Tuple[np.ndarray, np.ndarray]:
    """
    :param A: (n, n) symmetric real matrix
    :param eps:
    :return:
        w: (n,) eigenvalues sorted from largest to smallest
        v: (n, n) columns are normalized eigenvectors
        A = v diag(w) v^T
    """
    assert len(A.shape) == 2 and A.shape[0] == A.shape[1]
    assert np.isclose(A, A.T).all()

    T, V = lanczos_iteration(A)
    for T_t, U_t in qr_algorithm(T):
        if np.max(np.tril(T_t, k=-1)) < eps:
            w = np.diag(T_t)
            v = U_t @ V.T
            return w, v.T


def jacobi_eigh(S: np.ndarray, eps: float = 1e-3):
    assert len(S.shape) == 2 and S.shape[0] == S.shape[1]
    assert np.isclose(S, S.T).all()

    def givens_rotation(n: int, i: int, j: int, theta: float) -> np.ndarray:
        g = np.identity(n)
        c, s = np.cos(theta), np.sin(theta)
        g[i, i], g[j, j] = c, c
        g[i, j], g[j, i] = -s, +s
        return g

    def argmax_2d(A: np.ndarray) -> Tuple[int, int]:
        m, n = A.shape
        i = np.argmax(A)
        return i // n, i % n

    n = S.shape[0]

    v = np.identity(n)
    while True:
        # find pivot element
        S_tril_abs = np.abs(np.tril(S, k=-1))
        i, j = argmax_2d(S_tril_abs)
        if S_tril_abs[i, j] < eps:
            w = np.diag(S)
            return w, v.T
        # find givens matrix
        theta = np.arctan(2 * S[i, j] / (S[j, j] - S[i, i])) / 2
        g = givens_rotation(n, i, j, theta)
        # rotate
        S = g @ S @ g.T
        v = g @ v


def svd_from_eigh(A: np.ndarray, eps: float = 1e-3,
                  eigh: Callable[[np.ndarray, float], Tuple[np.ndarray, np.ndarray]] = jacobi_eigh) -> Tuple[
    np.ndarray, np.ndarray, np.ndarray]:
    """
    :param A:
    :param eps:
    :param eigh:
    :return:
    """
    assert len(A.shape) == 2
    m, n = A.shape
    if m > n:
        v, w, u = svd_from_eigh(A.T, eps=eps)
        return u, w, v

    w2_u, u = eigh(A @ A.T)
    rank = (w2_u > eps).sum()
    w = w2_u ** 0.5
    w_temp = w
    w_temp[rank:] = np.inf
    w_mo = w_temp ** (-1)
    v_temp = (np.diag(w_mo) @ u.T @ A).T[:, :rank]
    v = np.zeros(shape=(n, n))
    v[:, :rank] = v_temp
    for i in range(rank, n):
        v_i = np.random.normal(size=(n,))
        v_u, _ = gram_schmidt(v, v_i)
        v_i /= np.linalg.norm(v_i)
        v[i] = v_i
    return u, w, v


def conjugate_gradient(A: np.ndarray, b: np.ndarray) -> Iterable[np.ndarray]:
    """
    :param A: (n, n) symmetric, positive definite
    :param b: (n,)
    :return: x: (n,) solution of Ax = b
    """
    assert len(A.shape) == 2 and A.shape[0] == A.shape[1]
    assert np.isclose(A, A.T).all()
    assert sp.sparse.linalg.eigsh(A=A, k=1, which="SM")[0][0] > 0
    assert A.shape[0] == b.shape[0]

    x = np.zeros(shape=(n,))

    def residual(x: np.ndarray) -> np.ndarray:
        return b - A @ x

    D = np.zeros(shape=(n, n))

    for i in range(n):
        d = residual(x)
        d, _ = gram_schmidt(D, d, A)
        # for cg, lanczos, krylov basis -> only need to project to the last two directions
        d /= np.dot(d, A @ d) ** 0.5
        D[:, i] = d
        x += np.dot(d, b) * d
        yield x


if __name__ == "__main__":
    def take_n(i: Iterable[Any], n: int) -> List[Any]:
        return [x for _, x in zip(range(n), i)]


    n = 50
    m = 10
    A = np.random.normal(size=(n, n))
    b = np.random.normal(size=(n,))
    H = A.T @ A
    A_m = A[:m, :]
    H_m = A_m.T @ A_m

    #
    x = take_n(conjugate_gradient(H, b), n)[-1]
    print("conjugate gradient", np.max(np.abs(b - H @ x)))
    #
    Q, R = qr_decomposition(A)
    print("qr decomposition", np.max(np.abs(A - Q @ R)))
    #
    T, V = lanczos_iteration(H_m)
    print("lanczos iteration", np.max(np.abs(T - V.T @ H_m @ V)))
    #
    t = 100
    A_t, U_t = take_n(qr_algorithm(A), t)[-1]
    print("qr algorithm", np.max(np.abs(A_t - U_t @ A @ U_t.T)))
    #
    eps = 1e-6
    w, v = jacobi_eigh(H_m, eps=eps)
    print("jacobi eigen", np.max(np.abs(H_m - v @ np.diag(w) @ v.T)))
    #
    w, v = lanczos_qr_eigh(H_m, eps=eps)
    print("lanczos qr eigen", np.max(np.abs(H_m - v @ np.diag(w) @ v.T)))
    #
    u, w, v = svd_from_eigh(A_m, eps=eps)
    W = np.zeros(shape=(m, n))
    np.fill_diagonal(W, w)
    print("svd", np.max(np.abs(A_m - u @ W @ v.T)))
