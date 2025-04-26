# RLOG

Sequence paxos algorithm

# FEATURES

- [x] Consensus
- [x] Log Compaction
- [ ] Cluster membership change
- [ ] Persistent

# TODO 
- Cluster membership change
  - Maintain a version of ClusterState
  - Ignore all requests with different version (except Update)
  - ClusterState is changed from Decide (accepted from quorum)
  - Out of latest version nodes will be updated from Update