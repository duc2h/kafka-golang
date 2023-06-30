Play with kafka and golang
Something I can do:


* :x: Create a producer with wildcard topic.
    * The logic: 
        1. We have 2 topics: student_create, student_update.
        2. We have 3 consumers: 
            * A consumer consumes `student_create` topic.
            * A consumer consumes `student_update` topic.
            * A consumer consumes `student_create, student_update` topic.
* :heavy_check_mark: Create a group consumer with multiple topics. (3 consumers)
* :x: Deploy them like a cluster by k8s (PVC).
* :heavy_check_mark: Deploy kafka cluster (3 nodes)
    * Be careful with `host:port` in the cluster, read this [instruction](https://www.confluent.io/en-gb/blog/kafka-client-cannot-connect-to-broker-on-aws-on-docker-etc/#adding-new-listener)

* Some test:
    * :x: kafka TTL
    * :x: schema registry with protobuf: In order to see how schema-registry helpful .
    * :x: Stress test with k6.io
    * :heavy_check_mark: [Rebalancing in consumer group when a consumer join/leave happen](#rebalancing-in-consumer-group-when-a-consumer-joinleave-happen).
        *  [In order to add partition in runtime](#how-to-add-partitions-on-the-runtime)
    * :x: Offset, replay message when kafka down.
    * :x: Message key, order message follow partition.
        * Need to make a scenario for it.
    * :x: CDC with postgresql
    * :x: Implement out-box pattern.
    * ...



* Flow: post student -> publish message -> consume a message -> fetch user_id from user table -> Yes -> consume success remove message
                                                                                            -> No -> publish the message to dead letter, read the message proto,  


### Rebalancing in consumer group when a consumer join/leave happen
    * We create 4 partitions in `student_create` topic and 2 consumers in `group_consumer_student_create` consumer-group.
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

#### How to add partitions on the runtime
    1. Exec to kafka container: `docker exec -it [container-id] sh`
    2. Get help: `kafka-topics --help`
    3. Get describe of a topic: ` kafka-topics --describe --bootstrap-server [broker-server](eg: broker:9092) --topic [topic-name]`. Question: How about my kafka has 3 brokers?
    4. Add partition: `kafka-topics --bootstrap-server broker:9092 --topic [topic-name] --alter --partitions [number-of-partitions]`.
    5. Result
    ```
    sh-4.4$ kafka-topics --bootstrap-server broker1:9092 --topic student_create --describe
    Topic: student_create	TopicId: LODsLIlYQ9-21MM-pjIahw	PartitionCount: 4	ReplicationFactor: 3	Configs: 
        Topic: student_create	Partition: 0	Leader: 1	Replicas: 1,0,2	Isr: 1,0,2
        Topic: student_create	Partition: 1	Leader: 2	Replicas: 2,1,0	Isr: 2,1,0
        Topic: student_create	Partition: 2	Leader: 0	Replicas: 0,2,1	Isr: 0,2,1
        Topic: student_create	Partition: 3	Leader: 1	Replicas: 1,2,0	Isr: 1,2,0
    ```