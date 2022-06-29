# Order-service
***
The service provides functionality to work with orders.

# External requirements
***
    Go 1.18
    Docker
    Docker-compose
    Sqlc
    Goose
    Viper
    Zap
    GRPC


## Configuration

The service could be configured by providing environment variables.

| Name          | Meaning                                          | Example   |
|---------------|--------------------------------------------------|-----------|
| DBHOST        | Database connection host                         | 127.0.0.1 |
| DBPORT        | Database connection port                         | 5432      |
| DBUSER        | Database connection user                         | postgres  |
| DBNAME        | Database connection name                         | postgres  |
| DBPASSWORD    | Database connection password                     | password  |
| SSLMODE       | Database connection sslmode                      | disable   |
| RPCDRIVERPORT | Driver server RPC port                           | 5050      |
| RPCORDERPORT  | Order server RPC port                            | 5050      |
| RPCUSERPORT   | User server RPC port                             | 5050      |
| WAITINGTIME   | Maximum driver search time for user (in minutes) | 2         |
