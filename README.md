# BlockSight
Ethereum block insight / analytics 

+------------------+
|   External APIs  |
| (Infura/Alchemy) |
+--------+---------+
         |
         | Step 1: Retrieve blockchain data
         |
+--------v---------+
|  Data Ingestion  |
|     Pipeline     |
+--------+---------+
         |
         | Step 2: Publish data to Kafka topics
         |
+--------v---------+
|   Kafka Cluster  |
|                  |
|  Topics:         |
|  - Blocks        |
|  - Transactions  |
|  - Receipts      |
|  - Logs          |
+--------+---------+
         |
         | Step 3: Consume and process data
         |
+--------v---------+
|  Go Processing   |
|    Services      |
|                  |
|  - Block Parser  |
|  - Tx Processor  |
|  - Receipt Processor |
|  - Log Analyzer  |
+--------+---------+
         |
         | Step 4: Export metrics
         |
+--------v---------+
|   Prometheus     |
|     Metrics      |
|                  |
|  - Network Stats |
|  - Node Stats    |
|  - Contract Stats|
+--------+---------+
         |
         | Step 5: Visualize and monitor
         |
+--------v---------+
|     Grafana      |
|    Dashboard     |
|                  |
|  - Visualizations|
|  - Alerts        |
|  - Customization |
+------------------+