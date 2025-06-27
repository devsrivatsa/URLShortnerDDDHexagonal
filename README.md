# URL Shortener Service - DDD & Hexagonal Architecture

A high-performance URL shortening service built with Go, implementing Domain-Driven Design (DDD) and Hexagonal Architecture principles. The service provides fast and reliable URL shortening with support for multiple storage backends.

## Architecture Overview

### Domain-Driven Design Implementation

The application follows DDD principles with clear separation of concerns:

1. **Domain Layer** (`urlShortner` package)
   - Core domain model: `Redirect` entity
   - Domain interfaces: `RedirectService` and `RedirectRepository`
   - Business logic in the service implementation
   - Domain-specific errors and validation

2. **Application Layer** (`api` package)
   - HTTP handlers implementing the application's use cases
   - Request/Response handling
   - Content type negotiation (JSON/MessagePack)

3. **Infrastructure Layer**
   - Repository implementations (MongoDB/Redis)
   - Serialization adapters (JSON/MessagePack)
   - Configuration and environment management

### Hexagonal Architecture (Ports & Adapters)

The application implements hexagonal architecture through:

1. **Ports (Interfaces)**
   - Primary/Driving Ports: `RedirectHandler` interface
   - Secondary/Driven Ports: `RedirectRepository` and `RedirectSerializer` interfaces

2. **Adapters**
   - Primary Adapters: HTTP handlers in `api` package
   - Secondary Adapters:
     - Storage: MongoDB and Redis implementations
     - Serialization: JSON and MessagePack implementations

## System Design

### Components

1. **API Layer**
   - RESTful endpoints for URL shortening and redirection
   - Support for JSON and MessagePack serialization
   - Chi router with middleware for logging and recovery

2. **Storage Layer**
   - Pluggable storage backends
   - Supported databases:
     - Redis for high-performance caching
     - MongoDB for persistent storage
   - Easy to extend with new storage implementations

3. **Core Features**
   - URL shortening with unique code generation
   - Fast URL redirection
   - Validation of input URLs
   - Configurable storage backend

## API Endpoints

1. **Create Short URL**
   ```
   POST /
   Content-Type: application/json or application/x-msgpack
   
   Request Body:
   {
       "url": "https://example.com/long-url"
   }
   
   Response:
   {
       "code": "generated-short-code",
       "url": "https://example.com/long-url",
       "created_at": 1234567890
   }
   ```

2. **Redirect to Original URL**
   ```
   GET /{code}
   
   Response: 301 Redirect to original URL
   ```

## Getting Started

### Prerequisites

- Go 1.24.4 or later
- Docker and Docker Compose
- Redis or MongoDB (depending on your chosen storage backend)

### Environment Setup

1. Create a `.env` file in the root directory:
   ```env
   REPO_TYPE=redis  # or mongo
   REDIS_URL=redis:6379
   MONGO_URL=mongodb://mongodb:27017
   MONGO_DB=urlshortener
   MONGO_TIMEOUT=30
   ```

### Running with Docker Compose

1. Build and start the services:
   ```bash
   docker-compose up --build
   ```

This will start:
- The URL shortener service on port 8000
- Redis service on port 6379
- MongoDB service on port 27017 (if needed)

### Testing the Service

1. **Create a short URL (JSON)**:
   ```bash
   curl -X POST http://localhost:8000/ \
   -H "Content-Type: application/json" \
   -d '{"url": "https://example.com/very-long-url"}'
   ```

2. **Create a short URL (MessagePack)**:
   ```bash
   # Use the provided tool/test_msgpack.go for MessagePack testing
   go run tool/test_msgpack.go
   ```

3. **Access shortened URL**:
   ```bash
   curl -L http://localhost:8000/{generated-code}
   ```