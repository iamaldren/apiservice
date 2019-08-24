# apiservice

A simple API that exposes 5 endpoints that do simple transactions on a Redis server.

## Background

Ever since my friend and my colleague (2 different persons) introduced me to Golang, I've been quite drawn to it.

I tried reading tutorials, and watched some videos but wasn't able to finish it due to I'm really lazy to watch long videos, hence I just decided to create a very simple project where I can explore the language.

## Pre-requisites
- Golang
- Redis
- Docker (optional)

## Running the program

1. First you would need to have a Redis instance running or your machine.
   
   1.1 You can follow [here](https://redis.io/topics/quickstart), for instructions on how to install Redis.
   
   1.2 Or you can use Docker to host your Redis instance
   
   ```sh
     docker run --name <redis_name> -d redis:latest
   ```
   
   To connect via redis-cli, run the below
   
   ```sh
     docker run -it --network some-network --rm redis redis-cli -h some-redis
   ```
 2. Clone the project
 
 3. Open a command line, and go in the directory where the file main.go is placed.
    ```sh
      go run main.go
    ```
 4. The API services will be up and using port 8080.
 
 5. There will be 5 endpoints:
    * /v1/get/{key}
    */v1/set
    */v1/zset/{table}
    */v1/zgetall/{table}
    */v1/zget/{table}/{score}
    
## TODO

- Use dep for Go dependencies handling
- Add unit tests