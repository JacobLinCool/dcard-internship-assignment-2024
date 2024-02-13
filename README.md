# Dcard Internship Assignment 2024

Golang + MongoDB: 40K RPS API Server<sup>*</sup>

> <sup>*</sup> 40K RPS, 200% CPU usage in docker container on MBP M1 Pro 16GB, see [Third Iteration](#third-iteration)

## Assignment

The assignment can be found here: <https://drive.google.com/file/d/1dnDiBDen7FrzOAJdKZMDJg479IC77_zT/view>

TL;DR: Create a simple advertisement HTTP API server using Golang. External libraries and tools are allowed.

> It seems a quite fun project for me to refresh my Golang and MongoDB knowledge.

## Analysis

Numbers go first.

- ~3000 AD creations per day
- **~10000 AD views per second**
- ~1000 active ads at any given time

This assignment is quite simple. I don't want to over-engineer it with reverse proxies, load balancers, etc.

I think Golang is enough for this one, but considering the ease of use and scalability for data management, I choose MongoDB as the database.

> SQL or NoSQL is not a big deal for this assignment. Replication and sharding can always solve the problem I'm facing.

An RPS of 10k will be the primary challenge.

Golang is fast enough to handle this load, so the bottleneck will be the database.

Preliminary tests show MongoDB with indexes can handle around 2k RPS on my MBP, still far from the target.

Some cache mechanism is needed. Redis? No, only 1000 active ads, so I can cache them all in memory.

The tech stack becomes obvious:

- Golang / GIN for the HTTP server
- MongoDB for data storage
- Docker / Kubernetes for deployment (later)

## Design

Low latency/high throughput for the READ operation is the primary goal; no one wants to wait for an ad to load.
The WRITE operation is not that frequent, so it's not a big deal. Also, a few seconds of ad publishing delay is acceptable.

To reduce the latency, I will cache the ads in memory and update the local memory cache every second, so the cache will be eventually consistent.

```mermaid
graph TD
    P[Publisher]
    V[Viewer]
    S[Server]
    DB[Database]
    M[In-Memory Cache]
    P -->|Write| S -->|Write| DB -->|Update| M
    M -->|Cached| S -->|Query| V
```

## Result

I used K6 for load testing.

Environment:
- MacBook Pro M1 Pro 16GB, DevContainer (5 cores)
- Go 1.21.7
- MongoDB 7.0.2
- K6 0.49.0

### First Iteration

The result is quite good for the first iteration in DevContainer.

```bash
scenarios: (100.00%) 1 scenario, 300 max VUs, 40s max duration (incl. graceful stop):
        * default: 300 looping VUs for 10s (gracefulStop: 30s)


data_received..................: 193 MB 19 MB/s
data_sent......................: 10 MB  1.0 MB/s
http_req_blocked...............: avg=173.72µs min=375ns    med=792ns   max=175.38ms p(90)=1.87µs  p(95)=2.7µs   
http_req_connecting............: avg=164.91µs min=0s       med=0s      max=141.21ms p(90)=0s      p(95)=0s      
http_req_duration..............: avg=39.91ms  min=126.2µs  med=34.89ms max=214.87ms p(90)=66.72ms p(95)=78.79ms 
  { expected_response:true }...: avg=39.91ms  min=126.2µs  med=34.89ms max=214.87ms p(90)=66.72ms p(95)=78.79ms 
http_req_failed................: 0.00%  ✓ 0           ✗ 71399
http_req_receiving.............: avg=662.3µs  min=4.7µs    med=13.16µs max=157.87ms p(90)=70.83µs p(95)=163.75µs
http_req_sending...............: avg=105.1µs  min=2µs      med=4.2µs   max=137.76ms p(90)=11.54µs p(95)=28.54µs 
http_req_tls_handshaking.......: avg=0s       min=0s       med=0s      max=0s       p(90)=0s      p(95)=0s      
http_req_waiting...............: avg=39.14ms  min=104.66µs med=34.45ms max=183.84ms p(90)=65.37ms p(95)=76.49ms 
http_reqs......................: 71399  7135.616822/s
iteration_duration.............: avg=41.61ms  min=540.91µs med=35.99ms max=223.15ms p(90)=69.52ms p(95)=82.84ms 
iterations.....................: 71399  7135.616822/s
vus............................: 300    min=300       max=300
vus_max........................: 300    min=300       max=300
```

7K RPS, with less than 100% CPU usage for the server (K6 consumes the rest, 300%+, totaling 500% for 5 cores)
I believe the number can be doubled outside the DevContainer.

### Second Iteration

After some side work on k6 and disabling the debug log, the result is blazing fast compared to the first iteration.

> The k6 URL builder and logger IO caused the server cannot use all the CPU power.

```
scenarios: (100.00%) 1 scenario, 500 max VUs, 40s max duration (incl. graceful stop):
        * default: 500 looping VUs for 10s (gracefulStop: 30s)


data_received..................: 869 MB 87 MB/s
data_sent......................: 46 MB  4.6 MB/s
http_req_blocked...............: avg=15.99µs  min=250ns   med=417ns  max=41.18ms  p(90)=834ns   p(95)=1.16µs 
http_req_connecting............: avg=10.41µs  min=0s      med=0s     max=31.58ms  p(90)=0s      p(95)=0s     
http_req_duration..............: avg=15.03ms  min=49.95µs med=6.4ms  max=154.56ms p(90)=42.5ms  p(95)=54.04ms
  { expected_response:true }...: avg=15.03ms  min=49.95µs med=6.4ms  max=154.56ms p(90)=42.5ms  p(95)=54.04ms
http_req_failed................: 0.00%  ✓ 0            ✗ 323240
http_req_receiving.............: avg=844.43µs min=3.29µs  med=7.5µs  max=59.84ms  p(90)=58.7µs  p(95)=2.73ms 
http_req_sending...............: avg=36.36µs  min=1.37µs  med=2.29µs max=72.13ms  p(90)=4.29µs  p(95)=14.79µs
http_req_tls_handshaking.......: avg=0s       min=0s      med=0s     max=0s       p(90)=0s      p(95)=0s     
http_req_waiting...............: avg=14.15ms  min=43.12µs med=6.28ms max=140.83ms p(90)=39.69ms p(95)=48.96ms
http_reqs......................: 323240 32266.303532/s
iteration_duration.............: avg=15.41ms  min=65.37µs med=6.73ms max=154.58ms p(90)=43.03ms p(95)=54.67ms
iterations.....................: 323240 32266.303532/s
vus............................: 500    min=500        max=500 
vus_max........................: 500    min=500        max=500
```

32K RPS, with 250% CPU usage for the server (K6 consumes the rest, 200%+, totaling 500% for 5 cores)

> Note: VUS has been increased from 300 to 500.

### Third Iteration

After implementing some indexes for in-memory cache, 40K RPS is achieved.

```
scenarios: (100.00%) 1 scenario, 500 max VUs, 40s max duration (incl. graceful stop):
        * default: 500 looping VUs for 10s (gracefulStop: 30s)


data_received..................: 253 MB 25 MB/s
data_sent......................: 58 MB  5.8 MB/s
http_req_blocked...............: avg=21.11µs  min=250ns   med=416ns  max=41.4ms   p(90)=792ns   p(95)=1.12µs  
http_req_connecting............: avg=17.6µs   min=0s      med=0s     max=29.93ms  p(90)=0s      p(95)=0s      
http_req_duration..............: avg=11.82ms  min=29.54µs med=5.89ms max=124.07ms p(90)=30.98ms p(95)=39.29ms 
  { expected_response:true }...: avg=11.82ms  min=29.54µs med=5.89ms max=124.07ms p(90)=30.98ms p(95)=39.29ms 
http_req_failed................: 0.00%  ✓ 0            ✗ 406790
http_req_receiving.............: avg=529.52µs min=3.25µs  med=6.29µs max=48.1ms   p(90)=24.58µs p(95)=145.58µs
http_req_sending...............: avg=25.22µs  min=1.37µs  med=2.29µs max=51ms     p(90)=3.91µs  p(95)=12.58µs 
http_req_tls_handshaking.......: avg=0s       min=0s      med=0s     max=0s       p(90)=0s      p(95)=0s      
http_req_waiting...............: avg=11.27ms  min=23.33µs med=5.83ms max=117.43ms p(90)=29.39ms p(95)=35.63ms 
http_reqs......................: 406790 40651.191179/s
iteration_duration.............: avg=12.21ms  min=45.95µs med=6.31ms max=124.09ms p(90)=31.58ms p(95)=40.25ms 
iterations.....................: 406790 40651.191179/s
vus............................: 500    min=500        max=500 
vus_max........................: 500    min=500        max=500
```

40K RPS, with ~200% CPU usage for the server (K6 consumes the rest ~250%, totaling 500% for 5 cores)

### Other Observations

In addition to the previous results, I noticed that after restarting the Docker engine, the RPS increased to around 45K while maintaining the same CPU usage.

When running the application outside the DevContainer on my MBP, the RPS reached approximately 60K with around 300% CPU usage. The memory usage remained low, around 25MB when idle and 50MB under load.

Furthermore, I conducted another test by simulating 10,000 active ads, which is 10x the original number. This test was performed within the same DevContainer environment.

- The memory usage increased to around 60MB when idle and 100MB under load.
- The RPS achieved was approximately 17K with a CPU usage of around 300%.

```
scenarios: (100.00%) 1 scenario, 500 max VUs, 40s max duration (incl. graceful stop):
        * default: 500 looping VUs for 10s (gracefulStop: 30s)


data_received..................: 627 MB 63 MB/s
data_sent......................: 24 MB  2.4 MB/s
http_req_blocked...............: avg=16.26µs min=291ns  med=458ns  max=31.37ms  p(90)=792ns    p(95)=1.08µs  
http_req_connecting............: avg=12.05µs min=0s     med=0s     max=13.32ms  p(90)=0s       p(95)=0s      
http_req_duration..............: avg=28.94ms min=7.95µs med=5.87ms max=382.99ms p(90)=101.56ms p(95)=121.13ms
  { expected_response:true }...: avg=28.94ms min=7.95µs med=5.87ms max=382.99ms p(90)=101.56ms p(95)=121.13ms
http_req_failed................: 0.00%  ✓ 0            ✗ 172140
http_req_receiving.............: avg=99.62µs min=3.41µs med=9.29µs max=47.12ms  p(90)=23.75µs  p(95)=59.62µs 
http_req_sending...............: avg=12.92µs min=1.37µs med=2.25µs max=36.39ms  p(90)=3.7µs    p(95)=7.62µs  
http_req_tls_handshaking.......: avg=0s      min=0s     med=0s     max=0s       p(90)=0s       p(95)=0s      
http_req_waiting...............: avg=28.83ms min=0s     med=5.82ms max=382.98ms p(90)=101.4ms  p(95)=120.94ms
http_reqs......................: 172140 17183.435768/s
iteration_duration.............: avg=29.04ms min=51µs   med=5.92ms max=383.01ms p(90)=101.85ms p(95)=121.2ms 
iterations.....................: 172140 17183.435768/s
vus............................: 500    min=500        max=500 
vus_max........................: 500    min=500        max=500
```

## Conclusion

| Scenario          | RPS | CPU Usage | Memory Usage (idle / loaded) | Environment                            |
| ----------------- | --- | --------- | ---------------------------- | -------------------------------------- |
| 1,000 Active Ads  | 40K | ~200%     | 25MB / 50MB                  | DevContainer (5 cores), MBP M1 Pro 16G |
| 1,000 Active Ads  | 60K | ~300%     | -                            | MBP M1 Pro 16G                         |
| 10,000 Active Ads | 17K | ~300%     | 60MB / 100MB                 | DevContainer (5 cores), MBP M1 Pro 16G |

The architecture is designed to be simple, efficient, and easily scalable by adding additional server nodes. To handle a larger number of ADs, sharding by region (country) can be implemented.

```mermaid
graph TD
    DB[Database]
    S1["Server (TW)"]
    S2["Server (JP)"]
    S3["Server (US)"]
    M1["Local Cache (TW)"]
    M2["Local Cache (JP)"]
    M3["Local Cache (US)"]
    S1 -->|Write| DB -->|Update| M1
    S2 -->|Write| DB -->|Update| M2
    S3 -->|Write| DB -->|Update| M3
    M1 -->|Cached| S1
    M2 -->|Cached| S2
    M3 -->|Cached| S3
```

## Other stuff

- CI/CD has been set up with GitHub Actions, and the Docker image is pushed to GitHub Container Registry.
- Deployment is available in [`deployment`](./deployment) directory, including Docker Compose and Kubernetes.
