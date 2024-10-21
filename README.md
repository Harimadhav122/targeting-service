# delivery-service

    A delivery service is the microservice where an app/game is going to make a request. 
    Every request should contain the App ID, OS and Country information.
    Based on the rules in the targeting service, the delivery service will respond with a list of
    active campaigns to a particular user where the targeting rules are valid.

 ## Run

    ```go
    go run main.go
    ```

 ## Test

    Ensure mongodb server is running because the cache service is dependent on mongodb.
    Export env for mongodb connection uri and run the tests

    ```shell
    export MONGODB_CONN_URI="mongodb://localhost:27017/"
    ```

    ```go
    go test ./...
    ```

 ## Api

    http://localhost:8080/v1/delivery?app={app_id}&country={country_name}&os={os_name}

 ## HLA

 ![delivery-service-hla](https://github.com/user-attachments/assets/008c7aae-12a7-4326-b2c3-878dad97d3a2)
