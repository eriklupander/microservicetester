# microservicetester

### What is it?
Simplistic golang program that given a YAML spec boots a microservice landscape using docker-compose, uses Eureka endpoints to assert all services are up before performing some basic HTTP(s) REST endpoint tests. 

Supports OAuth tokens.

### Why should I use it?
Well, probably not :)

However, it may be a bit more easy to use compared to hacking your own bash/py/whatever script.

### How does it work?
You specify your setup in a .yml file:

    ---
    title: Microservices sample test file           # Give your setup a name!
    iterations: 1                                   # How many times should we loop this? Currently unsupported.
    host: http://192.168.99.100                     # Not used. Yet.
    docker_compose_root: /Users/myself/microservices-test      # Absolute path to the folder containing the docker-compose.yml file
    docker_compose_file: docker-compose.yml         # Name of your docker-compose.yml file. Well, it's typically this value.
    
    services:                                       # List of GET endpoints that must respond with HTTP status < 299 before the actual testing begins. This sample uses Eureka.
      - http://192.168.99.100:8761
      - http://192.168.99.100:8761/eureka/apps/config-server
      - http://192.168.99.100:8761/eureka/apps/edge-server
      - http://192.168.99.100:8761/eureka/apps/service-1
      - http://192.168.99.100:8761/eureka/apps/service-2
      - http://192.168.99.100:8761/eureka/apps/composite-service
    
    oauth:                                              # OAuth settings for getting a *drumroll* OAuth token.
      url: https://192.168.99.100:9999/uaa/oauth/token  # URL to the OAuth token provider service
      client_id: acme                                   # HTTP Basic auth password to call the OAuth service
      client_password: acmesecret                       # ditto password
      scope: webshop                                    # Some OAuth thing.
      grant_type: password                              # Another OAuth thing.
      username: user                                    # You actual username you want to have a Token issued for
      password: password                                # ditto password
      token_key: access_token                           # Yet another OAuth thing.
    
    endpoints:                                         # List of endpoints to test.
      - url: https://192.168.99.100/api/composite/composite-resource/123
        auth_method: TOKEN # TOKEN|BASIC|NONE
        method: GET
      
 
### Running it?
If you have the Go SDK installed, you can clone the repo and then run:

    go run src/github.com/eriklupander/mstest/*.go spec.yml
    
If you have a binay release:

    ./mstest spec.yml
    
    
The execution log might look like this after a successful run:

    Starting up...
    docker-compose installed OK
    Loaded specification 'Microservices sample test file'
    Docker starting up using /Users/myself/microservices-test/docker-compose.yml ...
    
    Waiting for all microservices to start...    
    
    http://192.168.99.100:8761                                 done                   
    http://192.168.99.100:8761/eureka/apps/config-server       done                   
    http://192.168.99.100:8761/eureka/apps/edge-server         done                   
    http://192.168.99.100:8761/eureka/apps/service-1           done                   
    http://192.168.99.100:8761/eureka/apps/service-2           done                   
    http://192.168.99.100:8761/eureka/apps/composite-service   done
    
    Getting OAuth token ... OK  
    
    Testing microservices...
    
    https://192.168.99.100/api/composite/composite-resource/123          ... OK                                          
    https://192.168.99.100/api/composite/composite-resource/999          ... OK                                          
    https://192.168.99.100/api/composite/composite-resource/888          ... OK                                          
    https://192.168.99.100/api/composite/composite-resource/777          ... OK                                                                                    
    
    All done.
    Docker shutting down...

    
# LICENSE

The MIT License (MIT)

Copyright (c) 2016 ErikL

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.