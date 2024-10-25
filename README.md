
# URL Shortener Project Documentation

This project is a URL Shortener API designed to handle high traffic and provide near real-time statistics. The project utilizes a modern tech stack, reactive programming, and cloud infrastructure to meet performance and scalability requirements.

---

## Table of Contents
1. [Project Overview](#project-overview)
2. [Tech Stack](#tech-stack)
3. [Architecture](#architecture)
4. [Setup Instructions](#setup-instructions)
5. [Configuration and Environment](#configuration-and-environment)
6. [Running Locally with Docker Compose](#running-locally-with-docker-compose)
7. [Deploying on GCloud EC2](#deploying-on-gcloud-ec2)

---

## Project Overview

This project implements a URL shortening service similar to goo.gl or bit.ly, with additional functionalities:
- Near real-time access statistics
- Ability to enable/disable or modify shortened URLs
- Support for high traffic, handling up to 1M RPM
- API REST with approximately 5000 RPM capacity, deployable on cloud infrastructure

## Tech Stack

- **Backend Language**: Go
- **Database**: MongoDB (public, no authentication)
- **Cache**: Redis (in-memory cache without authentication)
- **Cloud Infrastructure**: Google Cloud Platform (GCloud) and AWS EC2
- **Reactive Programming**: Using `rxgo` to handle non-blocking, concurrent operations
- **Containerization**: Docker & Docker Compose for local development

## Architecture

### Key Components:
1. **MongoDB**: Used to store URL information persistently. Exposed publicly on port `27017` without authentication (for testing purposes).
2. **Redis**: Acts as a cache layer to store frequently accessed URLs, improving response time and reducing database load.
3. **Reactive Programming**:
    - Using `rxgo` for non-blocking operations, such as creating and retrieving shortened URLs, caching, and managing statistics in real-time.
    - Reactive programming improves scalability and performance under heavy load.
4. **API Layer**: Provides RESTful endpoints for URL management, access, and real-time statistics.

### System Flow:
1. User accesses or creates a shortened URL.
2. Service processes the request reactively, updating Redis and MongoDB as necessary.
3. Real-time statistics are gathered to provide insight into API usage patterns.

---

## Setup Instructions

### Prerequisites
1. **Google Cloud SDK** and **AWS CLI** for managing cloud services.
2. **Docker** and **Docker Compose** for local development.
3. **Go** for backend development.
4. **rxgo** package for reactive programming (`go get -u github.com/reactivex/rxgo/v2`).

---

## Configuration and Environment

### MongoDB & Redis
- MongoDB is configured to run without authentication for this test environment.
- Redis is configured as a simple cache, also without authentication.

Both services are containerized with Docker Compose for easy setup.

### Environment Variables
- `MONGO_URI`: MongoDB URI, typically `mongodb://localhost:27017` for local testing.
- `REDIS_ADDRESS`: Redis URI, typically `localhost:6379`.

---

## Running Locally with Docker Compose

Use Docker Compose to set up MongoDB and Redis without authentication.

### `docker-compose.yml`
```yaml
version: '3.8'

services:
  mongo:
    image: mongo:latest
    container_name: mongo
    restart: always
    ports:
      - "27017:27017" # Expose MongoDB port publicly
    volumes:
      - mongo-data:/data/db
    command: mongod --bind_ip_all --noauth # No authentication, allow public access

  redis:
    image: redis:latest
    container_name: redis
    restart: always
    ports:
      - "6379:6379" # Expose Redis port without password
    volumes:
      - redis-data:/data

volumes:
  mongo-data:
    driver: local
  redis-data:
    driver: local
```

### Running Docker Compose
To start MongoDB and Redis, run:
```bash
docker-compose up -d
```

---

## Deploying on GCloud EC2

### Step-by-Step Guide

1. **Provision GCloud EC2 Instance**:
    - Create an EC2 instance on **Google Cloud** with sufficient resources to handle the anticipated load.
    - Ensure public access to MongoDB and Redis is configured securely using GCloud firewall rules.

2. **Set Up Environment**:
    - SSH into the instance and install **Docker** and **Docker Compose**.
    - Clone the project repository and navigate to the project directory:
      ```bash
      git clone https://github.com/yourusername/urlshortener.git
      cd urlshortener
      ```

3. **Deploy Docker Containers**:
    - Start MongoDB and Redis using Docker Compose:
      ```bash
      docker-compose up -d
      ```

4. **Run the API**:
    - Build and start the Go API:
      ```bash
      go build -o urlshortener
      ./urlshortener
      ```
    - Configure the Go API to connect to the MongoDB and Redis instances running on the GCloud instance.

5. **Expose API Publicly**:
    - Set up GCloud firewall rules to allow inbound HTTP/HTTPS requests to the instance.
    - You can use `ngrok` as an alternative during development for a public-facing endpoint without a complex setup.

6. **Autoscaling**:
    - Consider configuring **GCloud autoscaling** to add more instances under high traffic if running on a managed instance group.

### Example GCloud CLI Commands
For provisioning and firewall setup:
```bash
gcloud compute instances create urlshortener-instance     --zone=us-central1-a     --machine-type=e2-medium     --tags=http-server,https-server

gcloud compute firewall-rules create allow-http     --allow tcp:80

gcloud compute firewall-rules create allow-https     --allow tcp:443
```

---

## Future Enhancements

To make this application more robust and production-ready, consider the following improvements:

1. **Enhanced Security**:
    - Secure MongoDB and Redis with authentication, and restrict access to these services.
    - Use HTTPS for all public endpoints, especially for accessing sensitive APIs.

2. **Database Replication and Caching Strategies**:
    - Implement MongoDB replica sets to ensure high availability and failover capabilities.
    - Use a distributed caching system with Redis Cluster for scalability.

3. **Error Handling and Circuit Breaker Patterns**:
    - Integrate a circuit breaker pattern (e.g., with `go-resiliency`) to manage failures in external services like MongoDB or Redis.

4. **Rate Limiting and DDoS Protection**:
    - Implement rate limiting to prevent abuse of the API and secure the service from DDoS attacks.

5. **Comprehensive Testing**:
    - Expand testing to cover edge cases and stress test the API at higher RPMs to ensure resilience.

---

By following these steps, you will have a fully functional URL shortener service, capable of handling high loads, scalable with cloud infrastructure, and utilizing reactive programming for optimal performance. This document can serve as both a guide and a reference for maintaining and scaling the service in a real-world environment.

--- 
