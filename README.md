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
    * :heavy_check_mark: Rebalancing in consumer group when a consumer join/leave occur.
        1. 2 consumers in a group `group_consumer_student_create`.
        ```
            sh-4.4$ kafka-consumer-groups --bootstrap-server broker:9092 --group group_consumer_student_create --describe

            GROUP                         TOPIC           PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG             CONSUMER-ID                                 HOST            CLIENT-ID
            group_consumer_student_create student_update  0          2               2               0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
            group_consumer_student_create student_create  2          7               7               0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
            group_consumer_student_create student_create  3          8               8               0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
            group_consumer_student_create student_create  0          17              17              0               sarama-3ede9a42-8938-48ef-80e4-f4882d7dcd60 /172.18.0.1     sarama
            group_consumer_student_create student_create  1          25              25              0               sarama-3ede9a42-8938-48ef-80e4-f4882d7dcd60 /172.18.0.1     sarama
        ```

        2. delete a consumer in group `group_consumer_student_create`.
        ```
            sh-4.4$ kafka-consumer-groups --bootstrap-server broker:9092 --group group_consumer_student_create --describe

            GROUP                         TOPIC           PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG             CONSUMER-ID                                 HOST            CLIENT-ID
            group_consumer_student_create student_create  3          8               8               0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
            group_consumer_student_create student_update  0          2               2               0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
            group_consumer_student_create student_create  0          17              17              0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
            group_consumer_student_create student_create  2          7               7               0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
            group_consumer_student_create student_create  1          25              25              0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
        ```

        3. add new a consumer in group `group_consumer_student_create`.
        ```
            sh-4.4$ kafka-consumer-groups --bootstrap-server broker:9092 --group group_consumer_student_create --describe

            GROUP                         TOPIC           PARTITION  CURRENT-OFFSET  LOG-END-OFFSET  LAG             CONSUMER-ID                                 HOST            CLIENT-ID
            group_consumer_student_create student_create  0          17              17              0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
            group_consumer_student_create student_update  0          2               2               0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
            group_consumer_student_create student_create  1          25              25              0               sarama-ac9713b9-e465-48c4-9732-0400b9d45aa5 /172.18.0.1     sarama
            group_consumer_student_create student_create  2          7               7               0               sarama-c6675285-33d5-4605-b359-6481aa05c0ef /172.18.0.1     sarama
            group_consumer_student_create student_create  3          8               8               0               sarama-c6675285-33d5-4605-b359-6481aa05c0ef /172.18.0.1     sarama
        ```

    * :x: Offset, replay message when kafka down.
    * :x: Message key, order message follow partition.
        * Need to make a scenario for it.
    * :x: CDC with postgresql 
    * ...



post student -> publish message -> consume a message -> fetch user_id from user table -> Yes -> consume success remove message
                                                                   -> No -> publish the message to dead letter, read the message proto,  


* How to add partitions on the runtime.
    1. Exec to kafka container: `docker exec -it [container-id] sh`
    2. Get help: `kafka-topics --help`
    3. Get describe of a topic: ` kafka-topics --describe --bootstrap-server [broker-server](eg: broker:9092) --topic [topic-name]`. Question: How about my kafka has 3 brokers?
    4. Add partition: `kafka-topics --bootstrap-server broker:9092 --topic [topic-name] --alter --partitions [number-of-partitions]`.