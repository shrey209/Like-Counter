# Like-Counter

A scalable Like Counter system built in Go.

This service is designed to handle extremely high volumes of "like" events efficiently using a write-optimized architecture based on **Apache Kafka** and **Apache Cassandra**.

## ⚙️ Overview

The **Like-Counter** system is a backend service that ingests thousands of like events per second and performs optimized, batched writes to a database. It is ideal for applications such as social media platforms, content voting systems, or any system that requires fast and scalable like/increment operations.

### 🔧 Key Components

- **Language**: Golang (Go)
- **Database**: Apache Cassandra (chosen for its high write throughput and scalability)
- **Message Queue**: Apache Kafka (for high-throughput, asynchronous message handling)

## 🚀 How It Works

1. **Kafka Ingestion**  
   Like events are pushed into Kafka topics. Each message includes the content ID (e.g., `post_id`) that received the like.

2. **Batch Aggregation**  
   The service consumes messages in batches (e.g., thousands at a time). Instead of incrementing the database for each individual like, it aggregates counts in memory:

