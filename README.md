
# Gateway Service Project

Gateway service is a project that has the purpose of providing a gateway API for 2 other microservices, Users API & Find Number Position API.




## API Reference

The API uses basic auth authorization. This values have been set as enviroment variables. Here are the default ones that were set in the deployment:

Username: `admin`

Password: `password`

#### Default URL: http://localhost:8090

#### Get users

```http
  GET /users
```

#### Get item

```http
  GET /find_number_position?number={number}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `number`      | `int` | **Required**. Number of the position |

## Deployment

The project has a [docker-compose](https://github.com/docker/compose) configuration in order to deploy the service. Here's the command that you can use for deployment.

```bash
  make docker-compose
```

## Running Tests

To run tests, run the following commands

```bash
  make test
```

```bash
  make end-to-end-tests
```
## Environment Variables

Here are the environment variables per micro service that are defined


#### Gateway service


| Variable | Type     | Default value                |
| :-------- | :------- | :------------------------- |
| `PORT` | `int` | `8090` |
| `GRPC_CLIENT_KEEPALIVE_ALIVE_TIME` | `int` | `10` |
| `GRPC_CLIENT_KEEPALIVE_TIMEOUT` | `int` | `5` |
| `GRPC_CLIENT_PERMIT_WITHOUT_STREAM` | `boolean` | `false` |
| `GRPC_CLIENT_MAX_ATTEMPTS` | `int` | `5` |
| `GRPC_CLIENT_MAX_BACKOFF` | `string` | `0.01s` |
| `GRPC_CLIENT_BACKOFF_MULTIPLIER` | `float` | `1.0` |
| `GRPC_USERS_ADDRESS` | `string` | `localhost:50051` |
| `GRPC_FIND_NUMBER_POSITION_ADDRESS` | `string` | `localhost:50052` |
| `AUTH_USER` | `string` | `''` |
| `AUTH_PASSWORD` | `string` | `''` |

#### Users service


| Variable | Type     | Default value                |
| :-------- | :------- | :------------------------- |
| `GRPC_PORT` | `int` | `50051` |
| `REPOSITORY_FILE_DIRECTORY` | `string` | `/users_json` |
| `GRPC_SERVER_KEEPALIVE_ENFORCE_MIN_TIME` | `int` | `5` |
| `GRPC_SERVER_KEEPALIVE_ENFORCE_PERMIT_WITHOUT_STREAM` | `boolean` | `false` |
| `GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_IDLE` | `int` | `15` |
| `GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_AGE` | `int` | `30` |
| `GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_AGE_GRACE` | `int` | `5` |
| `GRPC_SERVER_KEEP_ALIVE_TIME` | `int` | `5` |
| `GRPC_SERVER_KEEP_ALIVE_TIMEOUT` | `int` | `1` |

#### Find Number Position service


| Variable | Type     | Default value                |
| :-------- | :------- | :------------------------- |
| `GRPC_PORT` | `int` | `50051` |
| `ARRAY_SIZE` | `int` | `100` |
| `GRPC_SERVER_KEEPALIVE_ENFORCE_MIN_TIME` | `int` | `5` |
| `GRPC_SERVER_KEEPALIVE_ENFORCE_PERMIT_WITHOUT_STREAM` | `boolean` | `false` |
| `GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_IDLE` | `int` | `15` |
| `GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_AGE` | `int` | `30` |
| `GRPC_SERVER_KEEPALIVE_MAX_CONNECTION_AGE_GRACE` | `int` | `5` |
| `GRPC_SERVER_KEEP_ALIVE_TIME` | `int` | `5` |
| `GRPC_SERVER_KEEP_ALIVE_TIMEOUT` | `int` | `1` |
