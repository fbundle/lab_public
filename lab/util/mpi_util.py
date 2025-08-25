from typing import Any, Iterable
from mpi4py import MPI
import threading
import queue

comm = MPI.COMM_WORLD
size = comm.Get_size()
rank = comm.Get_rank()


class Prog:
    def run(self):
        if rank == 0:
            self.ini_master()

            data_iter = iter(self.produce())
            q_emit = queue.Queue(maxsize=size - 1)  # queue of free workers
            q_send = queue.Queue(maxsize=size - 1)  # queue of messages to send
            q_recv = queue.Queue(maxsize=size - 1)  # queue of messages to recv
            for worker in range(1, size):
                q_emit.put(worker)

            def emit_task(q_emit: queue.Queue, q_send: queue.Queue, q_recv: queue.Queue):
                for data in data_iter:
                    worker = q_emit.get()
                    q_send.put((worker, data))

                q_send.put((worker, None))  # notify send_task that there is no more data to send

            def send_task(q_emit: queue.Queue, q_send: queue.Queue, q_recv: queue.Queue):
                while True:
                    worker, data = q_send.get()
                    if data is None:  # no more data to send, break
                        q_recv.put(False)  # notify recv_task that there is no more data to recv
                        break
                    comm.send(data, dest=worker)
                    q_recv.put(True)

            def recv_task(q_emit: queue.Queue, q_send: queue.Queue, q_recv: queue.Queue):
                while True:
                    token = q_recv.get()
                    if not token:  # no more data to recv, break
                        break
                    worker, data_out = comm.recv()
                    self.consume(data_out)
                    q_emit.put(worker)

            t_list = [
                threading.Thread(target=target, args=(q_emit, q_send, q_recv))
                for target in [emit_task, send_task, recv_task]
            ]
            for t in t_list:
                t.start()
            for t in t_list:
                t.join()
            # tell all workers to stop
            for worker in range(1, size):
                comm.send(None, dest=worker)

            self.del_master()

        else:
            # worker
            self.ini_worker()
            while True:
                data = comm.recv()
                if data is None:
                    break
                data_out = self.apply(data)
                comm.send((rank, data_out), dest=0)
            self.del_worker()

    def ini_master(self):
        raise NotImplemented

    def del_master(self):
        raise NotImplemented

    def ini_worker(self):
        raise NotImplemented

    def del_worker(self):
        raise NotImplemented

    def apply(self, item: Any) -> Any:
        raise NotImplemented

    def produce(self) -> Iterable[Any]:
        raise NotImplemented

    def consume(self, result: Any):
        raise NotImplemented
