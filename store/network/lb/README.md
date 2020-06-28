### Load Balancer partitioning

This section implements partitioning strategies, at a network level.
Every corresponding router implementation, routes traffic to the specific nodes of the distributed store.
- random
- sharded
- hashing
- leader-follower