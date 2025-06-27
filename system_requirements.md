## Requirements

#### Functional Requirements

1. **URL Shortening**

   - Generate unique shortened URLs for long URLs
   - Ensure uniqueness and consistency of shortened URLs

2. **URL Redirection**

   - Fast and efficient redirection from short URL to original URL
   - Support for HTTP 301/302 redirects

#### Non-Functional Requirements

1. **Performance**
   - Minimize redirect latency
   - Optimize for read-heavy operations
   - Efficient caching strategy

## Capacity Planning

### Traffic Estimates

- 100M daily active users
- 1B reads per day
- 10K requests per second
- 1-5 total lifetime URLs to store

### Write Operations

- ~1000 URLs generated per second
- 31.5B URLs created per year
- System designed to last 10 years with 7^62 possible URLs

### Read Operations

- 10:1 read to write ratio
- ~300B reads per year

## Technical Specifications

### Storage Requirements

- Record size: ~1KB
- Total storage for 1B URLs: ~5TB
- Multiple database instances to handle 10K req/sec
