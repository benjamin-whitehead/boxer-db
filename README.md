# boxer-db
BoxerDB is a toy distributed key value store implemented in Go. BoxerDB uses leader based replication. There are still issues with it that need to be fixed, and it is not intended to be used in any situation. It is just made so I could learn more about Distributed Systems and Databases.

 ## Features:
 - A key value store that allows for GET, PUT and DELETE operations via an HTTP API.
 - Replicates all operations on leader to follower nodes.
 - Persists all data to a file on disk, and reads saved state on startup.

 ## Usage:
 To use the database, send requests to the HTTP API. The database currently supports the following requests:

 ### To create a entry:
```
    POST http://localhost:8080/api/v1/<key>
    {
        "value": "Hello World!"
    }
 ```
Note: since the type used for the values are interfaces, you are able to store JSON values in the database. Example:
```
    POST http://localhost:8080/api/v1/<key>
    {
        "value": {
            "message": "JSON!",
            "time": 100
        }
    }
 ```
 This endpoint also supports the use of PUT.

 Successful response
 ```
 200 OK
 ```

 Unsuccessful responses:
 ```
 403 FORBIDDEN
{
	"message": "This action is only available to the leader"
}
``` 
This occurs when you try and make a request to the replica nodes, as the replica is only used for consistency, and is read only.

```
404 NOT FOUND
```
This occurs when the key is not found in the database.

 ### To read an entry:
 ```
    GET http://localhost:8080/api/v1/<key>
 ```
 Successful response:
 ```
    200 OK
    {
        "value": "Hello World!"
    }
```
Unsuccessful response:
```
404 NOT FOUND
```

### To delete an entry:

```
    DELETE http://localhost:8080/api/v1/<key>
 ```
 Successful response
 ```
 200 OK
 ```
 Unsuccessful responses:
 ```
 403 FORBIDDEN
{
	"message": "This action is only available to the leader"
}
``` 
This occurs when you try and make a request to the replica nodes, as the replica is only used for consistency, and is read only.
```
404 NOT FOUND
```
This occurs when the key is not found in the database.


 ## Installation:
 This is designed to work with Docker and Docker Compose. Since it is not designed to work on actual hardware, getting it to run on different machines can be done by cloning the repository, building the binary, and running it on multiple machines.

 The easier way of running it, is to use Docker Compose. Simply clone the repository, and run:
 ```
 docker-compose up --build
 ```
Which will start a local cluster with one leader and three followers, that you will be able to make requests to on localhost:8080. The configuration for each node is defined in the docker compose file.

## Issues:
There are some issues that still need to be worked out. One issue is with the Docker Compose file. I am still getting used to Docker and Docker Compose, and there is an issue where if you stop the docker compose cluster, the data that is saved into the file is lost. This is because I am having issues with the volumes.


 ## Future Work:
 In the future, I want to implement the following features:
  - Integration with configuration manager (While this is designed to be consistent, the system is uncapable of recovering from any node failures. I am currently working on a leaderless configuration manager which will allow for leader elections and nodes joining the system. This is the highest priority for me.)
 - Implement async clients (handle multiple requests concurrently)
 - Expand data model (Support collections and multiple databases)
