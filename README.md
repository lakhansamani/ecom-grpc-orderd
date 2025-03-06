# orderd

Order service for ecom-grpc example

```
docker run --name postgres-cluster -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
docker exec -it postgres-cluster psql -U postgres -c "CREATE DATABASE orderdb;"
```

```sh
grpcurl -plaintext -H "authorization: bearer JWT_TOKEN"  -d '{ "product": "book 1", "quantity": 1  }' -proto=apis/order/v1/order.proto localhost:50052 order.v1.OrderService/CreateOrder

grpcurl -plaintext -H "authorization: bearer JWT_TOKEN"  -d '{ "id": "ID"  }' -proto=apis/order/v1/order.proto localhost:50052 order.v1.OrderService/GetOrder
```