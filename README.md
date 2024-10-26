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
    <img width="382" alt="Screenshot 2024-10-26 at 2 49 23 PM" src="https://github.com/user-attachments/assets/edab096d-0c80-4792-92a1-3a38348011ea">

 ## Api

    http://localhost:8080/v1/delivery?app={app_id}&country={country_name}&os={os_name}

 ## HLA
    ![delivery-service-hla](https://github.com/user-attachments/assets/a84dc5ea-56e6-4198-9304-26876511aeba)
