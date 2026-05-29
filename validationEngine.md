VALIDATION RULES :: 



IF: High traffic system (social media, streaming, etc.)
AND: No cache layer between App Server and Database
THEN: ⚠️ Warning — Every request hits the database directly
WHY: Database calls are 100x slower than cache reads
      At scale this causes DB overload and slow response times
SUGGESTION: Add Redis or Memcached between your app server and DB


IF: File uploads in your system
AND: Files going directly to Database
THEN: ⚠️ Wrong tool for the job
WHY: Databases aren't built for binary file storage
SUGGESTION: Use Object Storage (S3/Cloudflare R2) for files


IF: High write traffic expected
AND: No message queue present
AND: Direct connection from Backend → Database
THEN: ⚠️ DB will get overwhelmed during traffic spikes
WHY: No buffer between your backend and DB
SUGGESTION: Add Kafka for high volume or 
            RabbitMQ for task processing


IF: MongoDB selected
AND: User describes relationship-heavy features 
     (followers, friends, recommendations)
THEN: ⚠️ Wrong database for this use case
WHY: Document databases handle relationships 
     poorly at scale
SUGGESTION: Consider PostgreSQL for relational data


IF: Multiple operations triggered by single user action
AND: All operations are synchronous
AND: No message queue present
THEN: ⚠️ High response time and potential timeout
WHY: User waits for every operation to complete
SUGGESTION: Use queue for non-critical async operations
            Only keep critical path synchronous


IF: Multiple operations triggered by single user action
AND: All operations are synchronous
AND: No message queue present
THEN: ⚠️ High response time and potential timeout
WHY: User waits for every operation to complete
SUGGESTION: Use queue for non-critical async operations
            Only keep critical path synchronous


IF: Frontend directly connected to Database
AND: No backend layer in between
THEN: 🔒 Critical Security Issue
WHY: DB credentials exposed in browser
     Anyone can access your entire database
SUGGESTION: Always put a backend between 
            frontend and database


IF: Backend API present
AND: No authentication layer detected
THEN: 🔒 Security Issue
WHY: Anyone can call your API endpoints
     No way to verify who is making requests
SUGGESTION: Add JWT / OAuth2 / Session based auth


IF: Authentication present
AND: No authorization checks
THEN: 🔒 Security Issue
WHY: Logged in users can access other users data
     No role based restrictions
SUGGESTION: Implement role based access control (RBAC)


IF: Backend connected to SQL Database (PostgreSQL)
AND: No input validation layer detected
THEN: 🔒 Security Issue
WHY: Raw user input reaching database queries
     Attackers can manipulate or destroy your DB
SUGGESTION: Use parameterized queries
            Add input validation middleware


IF: Backend present
AND: No error handling layer detected
THEN: 🔒 Security Issue
WHY: Stack traces expose system internals
     Attackers use this for reconnaissance
SUGGESTION: Add global error handler
            Never expose stack traces in production
            Use internal logging instead


IF: Authentication endpoint present
AND: No rate limiting detected
THEN: 🔒 Security Issue
WHY: Unlimited login attempts allow brute force
     Attacker can try millions of passwords
SUGGESTION: Add rate limiting to auth endpoints
            Lock account after N failed attempts
            Consider captcha for suspicious activity


IF: System has frontend and backend
AND: No HTTPS/TLS layer detected
THEN: 🔒 Security Issue
WHY: Data travels as plain text over network
     Anyone on same network can read it
SUGGESTION: Always use HTTPS in production
            Get SSL certificate (free via Let's Encrypt)
            Redirect all HTTP to HTTPS


IF: Single backend server
AND: No load balancer present
AND: High traffic expected
THEN: 🔧 Reliability Issue
WHY: Single point of failure
     One crash = entire app goes down
SUGGESTION: Add load balancer
            Run minimum 2 server instances
            Implement health checks


IF: Single database instance
AND: No replica detected
AND: Data loss unacceptable
THEN: 🔧 Reliability Issue
WHY: Database is single point of failure
     Crash = app completely dead
     No failover available
SUGGESTION: Set up primary replica replication
            Configure automatic failover
            Use synchronous replication for 
            critical data


Without redundancy:
Users → Single Server → Single Database
         ↓ crashes          ↓ crashes
      app dead           data unreachable

With redundancy:
Users → Load Balancer → Server 1 ✅
                      → Server 2 ✅
                      → Server 3 ✅
                           ↓
                      Primary DB ✅
                           ↓ replicates
                      Replica DB ✅

               
IF: Database present
AND: No backup strategy detected
THEN: 🔧 Reliability Issue
WHY: Human errors, bugs, or attacks can 
     cause permanent data loss
     No way to recover without backups
SUGGESTION: Enable automated daily backups
            Implement Point in Time Recovery
            Test your backups regularly
            Store backups in separate location


IF: Backend and database present
AND: No monitoring layer detected
THEN: 🔧 Reliability Issue
WHY: You find out about failures from 
     angry users not from your system
     No visibility into performance degradation
SUGGESTION: Add monitoring and alerting
            Monitor infrastructure + performance
            + business metrics
            Use tools like Prometheus + Grafana
            or Datadog


IF: Backend servers present
AND: No deployment strategy defined
THEN: 🔧 Reliability Issue
WHY: Direct deployments risk full downtime
     Bugs affect all users immediately
     Rollback is slow and painful
SUGGESTION: Implement Blue-Green for zero 
            downtime deployments
            Or Canary for gradual rollouts
            Combined with CI/CD pipeline



            IF: Backend making sequential DB calls
AND: Data could be fetched in parallel
THEN: ⚡ Performance Issue
WHY: Sequential calls add up response time
     N+1 queries can make 1000+ DB calls
     For a single page load
SUGGESTION: Use parallel queries where possible
            Use JOIN queries instead of 
            multiple sequential calls
            Implement eager loading


IF: Backend making sequential DB calls
AND: Data could be fetched in parallel
THEN: ⚡ Performance Issue
WHY: Sequential calls add up response time
     N+1 queries can make 1000+ DB calls
     For a single page load
SUGGESTION: Use parallel queries where possible
            Use JOIN queries instead of 
            multiple sequential calls
            Implement eager loading


IF: Database present
AND: Large dataset expected
AND: No indexes defined on search columns
THEN: ⚡ Performance Issue
WHY: Full table scans on large datasets
     Can take seconds per query
     Gets worse as data grows
SUGGESTION: Add indexes on frequently 
            searched columns
            Be selective — don't index everything
            Avoid indexes on write-heavy columns


IF: Frontend present
AND: App serves static assets
AND: No CDN detected
AND: Global or large user base expected
THEN: ⚡ Performance Issue
WHY: Static files served from single origin
     High latency for distant users
     Origin server handles unnecessary load
SUGGESTION: Add CDN for static assets
            Use Cloudflare, AWS CloudFront,
            or similar
            Cache static assets at edge


IF: Backend performs heavy computation
AND: Same computation repeated per request
AND: No caching strategy detected
THEN: ⚡ Performance Issue
WHY: Redundant computation wastes resources
     Response time increases with data size
     DB gets hammered unnecessarily
SUGGESTION: Cache computed results
            Use TTL for rarely changing data
            Use cache invalidation for 
            user specific data
            Pre compute for analytics/aggregates



IF: Multiple backend servers present
AND: User authentication exists
AND: No centralized session storage
THEN: 📈 Scalability Issue
WHY: Sessions stored on individual servers
     Users get logged out randomly
     As requests hit different servers
SUGGESTION: Use Redis for centralized 
            session storage
            All servers share same session store
            Stateless servers scale freely


IF: Backend servers present
AND: Server stores local state
AND: Horizontal scaling intended
THEN: 📈 Scalability Issue
WHY: Stateful servers can't be 
     scaled horizontally
     Requests must hit same server
     Defeats purpose of multiple servers
SUGGESTION: Make servers stateless
            Move all state to Redis or DB
            Use JWT instead of server sessions


IF: Single monolith architecture
AND: Large team size indicated
AND: High scale expected
AND: Multiple distinct domains present
THEN: 📈 Scalability Issue
WHY: Monolith becomes bottleneck at scale
     Can't scale individual components
     Team velocity slows down
     Single deployment risk
SUGGESTION: Consider microservices
            Split by business domain
            Each service owns its data
            BUT: Don't over-engineer early stage


IF: Microservices architecture
AND: Small team indicated (< 10)
AND: Early stage product
THEN: 📈 Scalability Warning
WHY: Microservices add operational complexity
     Small teams get overwhelmed managing
     Multiple services, deployments, networks
SUGGESTION: Start with monolith
            Extract services as you grow
            Premature microservices = 
            distributed monolith


IF: Microservices architecture
AND: Services call each other synchronously
AND: No circuit breaker pattern detected
THEN: 📈 Scalability Issue
WHY: One failing service cascades
     Takes down dependent services
     Entire system fails from one point
SUGGESTION: Implement circuit breaker pattern
            Define fallback responses
            Use graceful degradation
            Tools: Hystrix, Resilience4j


IF: Cloud infrastructure indicated
AND: Variable traffic expected
AND: No auto scaling configured
THEN: 📈 Scalability Issue
WHY: Manual scaling is slow and reactive
     Over provisioning wastes money
     Under provisioning affects users
SUGGESTION: Configure auto scaling rules
            Set CPU and memory thresholds
            Use predictive scaling for 
            known traffic patterns
            Minimum 2 instances always running


Rule 028:
IF: Redis selected
AND: Used only for basic caching
AND: No sessions or queues
THEN: 💡 Suggestion
      Memcached might be simpler
      and faster for pure caching


Rule 029:
IF: Memcached selected
AND: Sessions need to be stored
THEN: ⚠️ Warning
      Memcached has no persistence
      Sessions will be lost on restart
      Use Redis instead

      