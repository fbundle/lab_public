import numpy as np

def block_krylov(A: np.ndarray, k: int, q: int = 5, count_sketch: bool = False, cheb: bool = False) -> tuple[np.ndarray, np.ndarray]:
    """
    :param A: (n, d) matrix
    :param k: rank
    :param q: number of blocks
    :param count_sketch: whether to use count_sketch matrix instead of gaussian
    :return: (Z, P) where Z \approx U_k, P \arppox \Sigma_k V_k^T, U_k \Sigma_k V_k^T \approx A
    """
    assert len(A.shape) == 2
    n, d = A.shape
    assert k < min(n, d)

    if count_sketch:
        h = np.random.randint(low=0, high=d, size=(k,))
        phi = np.zeros(shape=(d, k))
        for i in range(k):
            phi[h[i], i] = 1
        d = np.random.choice([-1, +1], size=(k, ))
        Pi = phi @ np.diag(d)
        
    else:
        Pi = np.random.normal(loc=0, scale=1, size=(d, k))
    
    def column_normalize(u: np.ndarray) -> np.ndarray:
        return u / np.linalg.norm(u, axis=0)

    AAT = A @ A.T
    K_list = [column_normalize(A @ Pi)]
    if cheb:
        K_list.append(column_normalize(2 * AAT @ K_list[-1]))
        for i in range(2, q):
            K_list.append(column_normalize(2 * AAT @ K_list[-1] - K_list[-2]))
    else:
        for i in range(1, q):
            K_list.append(column_normalize(AAT @ K_list[-1]))
        

    K = np.concatenate(K_list, axis=1)

    Q, R = np.linalg.qr(K, mode="reduced")

    U, S, Vh = np.linalg.svd(Q.T @ A)
    U_k = U[:, 0:k]

    Z = Q @ U_k
    P = Z.T @ A

    return Z, P


if __name__ == "__main__":
    a = np.random.normal(size=(2000, 5000))
    k = 100

    u, s, vh = np.linalg.svd(a)
    z = u[:, 0:k]

    z_1, p_1 = block_krylov(a, k)
    print(np.abs(z - z_1).mean())