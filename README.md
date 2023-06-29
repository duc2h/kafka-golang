Play with kafka and golang
Something I can do:


* :x: Create a producer with wildcard topic.
    * The logic: 
        1. We have 2 topics: student_create, student_update.
        2. We have 3 consumers: 
            * consumer consumes `student_create` topic.
            * consumer consumes `student_update` topic.
            * consumer consumes `student_*` topic.
* :heavy_check_mark: Create a group consumer with multiple topics. (2 consumers)
* :x: Deploy them like a cluster by k8s (PVC).
* Some test:
    * :x: kafka TTL
    * :x: schema registry with protobuf: In order to see how schema-registry helpful .
    * :x: Stress test with k6.io
    * :x: Rebalancing in consumer group when a consumer join/leave occur.
    * :x: Offset, replay message when kafka down.
    * :x: Message key, order message follow partition.
    * :x: CDC with postgresql 
    * ...



post student -> publish message -> consume a message -> fetch user_id from user table -> Yes -> consume success remove message
                                                                   -> No -> publish the message to dead letter, read the message proto,  


* How to add partitions on the runtime.
    1. Exec to kafka container: `docker exec -it [container-id] sh`
    2. Get help: `kafka-topics --help`
    3. Get describe of a topic: ` kafka-topics --describe --bootstrap-server [broker-server](eg: broker:9092) --topic [topic-name]`. Question: How about my kafka has 3 brokers?
    4. Add partition: `kafka-topics --bootstrap-server broker:9092 --topic [topic-name] --alter --partitions [number-of-partitions]`.