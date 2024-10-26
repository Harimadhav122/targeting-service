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

    ```go
    go test ./... -coverprofile cover.out
    go tool cover -html=cover.out
    ```

 ## Coverage
    
 ## Api

    http://localhost:8080/v1/delivery?app={app_id}&country={country_name}&os={os_name}

 ## HLA

 ![delivery-service-hla](https://github.com/user-attachments/assets/008c7aae-12a7-4326-b2c3-878dad97d3a2)
