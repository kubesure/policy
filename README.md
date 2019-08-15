# policy

protoc --proto_path=..\api\v1\proto --go_out=plugins=grpc:vendor/github.com/kubesure/policy/api/v1 ..\api\v1\proto\publisher.proto

### Setup

   ```
      CREATE DATABASE `policy`;
      USE policy;
      DROP TABLE policy;
      CREATE TABLE `policy` (
         `id` int(11) NOT NULL AUTO_INCREMENT,
         `quote_no` int(16) DEFAULT NULL,
         `receipt_no` int(10) DEFAULT NULL,
         `created_date` timestamp DEFAULT now(),
         `updated_date`  timestamp DEFAULT now(),
         PRIMARY KEY (`id`)
      );
    ```

### Run and Test
```
 go run policy.go   
```
curl -X POST http://localhost:8000/api/v1/health/poicies -H 'Content-Type: application/json' -d '{
    "QuoteNumber": "12343456",        
    "ReceiptNumber": "1234345678"     
} '
```
