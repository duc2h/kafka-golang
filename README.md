Play with kafka and golang
Something I can do:


* Create a producer with wildcard topic.
* Create a group consumer with multiple topics. (3 consumers)
* Deploy them like a cluster by k8s (PVC).
* Some test:
    * kafka TTL
    * schema registry with protobuf: In order to see how schema-registry helpful .
    * Stress test with k6.io
    * Rebalancing in consumer group when a consumer join/leave occur.
    * Offset, replay message when kafka down.
    * Message key, order message follow partition.
    * CDC with postgresql 
    * ...
