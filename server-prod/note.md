# Understanding HTTP Server Configuration in GO

### How server request response cycle happen
Browser => HTTP Request => HTTP Server => Go code => HTTP Response

## Why do we need configuration?
- Removes code repetition, keeps all server configuration in one place
- Ease to change environment, example test config could be different form production or development configs
- go has a default value of 0 timeout values which pose a serious security vulnerability flagged as a CWE-770 (unrestricted Resource allocation) vulnerability

#### Host (string)
"0.0.0.0" is used in prod server so that anyone on the internet can access your service

### Port (int)
(eg 443) Ports are number and they act as a door into the application

### ReadTimeout (time.Duration)
more like saying if the server takes too long to send the request, close the connection. Example imagine 5000 customers in your restuarant not buying but slowly reading the menu, this waste resources and you dont make any money.  the recommended value is 10 - 30 s

we use time.Duration because it lets us write the values like this;

15 * times.Second

instead of 15

you can agree that the former is more easier to read.

### WriteTimeout (time.Duration)
same thing with the read only that the client receives its request rather slowly. the recommended value is 10 - 30 s

### IdleTimeout (time.Duration)
Keeping connections open is good because it allows us reuse those connections time and time again, but when nobody is using the connection should we keep it alive forever NO, we ought to close is, so we set a maxtime for which if not reused we close it. the recommended valur is 60 - 120 s

### ReadHeaderTimeout (time.Duration)
attackers might want to sent the header extremely slowly this nature of attack is slowloris attack.  the recommended valur is 5 - 10 s

### MaxHeaderBytes (int)
This is usually small and prevents people from maxing resource usage due to heavy data

### MaxBodySize (int64)
usually larger than the header this is used to prevent resource exhaustion as well. You should configure this based on the solution your service is rendering.

### RateLimit (float64) 
this is used to prevent one user from blocking the app for others to use, by sending multiple request within a very short window.

### RateLimitBurst (int)
Describes how mnay more extra requests can arrive simultaneously after being blocked. Example there is a RateLimit of 1req/m and we have a burst of 5req/m it means that if 5 people come into the cafe together we can accomdate them and kick out the next customer coming since the 

### EnableTLS (bool)
if this is true packets would be encrypted during transaction, else anybody who intercepts the network would read it.

### CertFile (string)
HTTPS server need cert. this file proves the server's identity

### KeyFile (string)
this is used with the cert to establish secure connections. it must be protected carefully and never shared publicly

## Rate Limiting in Go - Protecting your server from overload

we have specified an interface for the rate limiter of choice this is a matter of flexibility as it enables us change our preferred ratelimiter of choice in the future if we need to so tomorrow we have write implementation for uber or redis based rate limiter if we want to and it must support Allow behaviour.

Rate Limiter is the first line of defense for abusing your service, not the final bustop. It is majorly concerned about solving problem centered around server resource mgt.

Don't keep the burst small a single resource can in turn make several other requests simultaneously

Dont' use one global limiter, use a limiter per user, API key, or IP address. So one user do not block the process for everyone.

## Database connection pooling in Go
*The concept:*

When you make request to the database. 

It reaches for the database, authenticates the requester, established a connection and make the resources available to the requester.

Opening a connection each time we make a request is expensive and wastes resources, so we keep a number of connections in the connection pools open and each time we make a request a connection is borrowed, command executed and when it is done returns the connection to the pool. The connection is reused, not recreated.

The MaxIdleConns is best to equal the MaxOpenConns, however know that an increased number of connections doesn't mean faster application it means more money as a result of its consumed resource usage
