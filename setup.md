### policy

protoc --proto_path=..\api\v1\proto --go_out=plugins=grpc:vendor/github.com/kubesure/policy/api/v1 ..\api\v1\proto\publisher.proto

### Setup Dev

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

### Setup K8s cluster

Install percona mysql 

```
   git clone https://github.com/kubesure/helm-charts
   cd ./helm-charts/percona-xtradb-cluster  
   helm install . --name=mysql-policy
```

Create db

Follow [notes](#Helm install notes mysql policy db) for mysql login credentials

```
   kubectl exec --namespace default -ti mysql-policy-pxc-0 -c database -- mysql -uroot -prootpassword
```   

```
   use policy;
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

```
   curl -i -X POST http://localhost:8000/api/v1/health/poicies -H 'Content-Type: application/json' -d '{
      "QuoteNumber": "12343456",        
      "ReceiptNumber": "1234345678"     
   } '
```

### Helm install notes mysql policy db

NOTES:
Percona can be accessed via port 3306 on the following DNS name from within your cluster:
mysql-policy-pxc.default.svc.cluster.local

To get your root password run (this password can only be used from inside the container):

    $ kubectl get secret --namespace default mysql-policy-pxc -o jsonpath="{.data.mysql-root-password}" | base64 --decode; echo

To get your xtradb backup password run:

    $ kubectl get secret --namespace default mysql-policy-pxc -o jsonpath="{.data.xtrabackup-password}" | base64 --decode; echo

To check the size of the xtradb cluster:

    $ kubectl exec --namespace default -ti mysql-policy-pxc-0 -c database -- mysql -e "SHOW GLOBAL STATUS LIKE 'wsrep_cluster_size'"

To connect to your database:

1. Run a command in the first pod in the StatefulSet:

    $ kubectl exec --namespace default -ti mysql-policy-pxc-0 -c database -- mysql

2. Run a percona pod that you can use as a client:

    $ kubectl run -i --tty --rm percona-client --image=percona:5.7.19 --restart=Never -- mysql -h mysql-policy-pxc.default.svc.cluster.local -upolicy \
      -p$(kubectl get secret --namespace default mysql-policy-pxc -o jsonpath="{.data.mysql-password}" | base64 --decode; echo) \
     policy

To view your Percona XtraDB Cluster logs run:

  $ kubectl logs -f mysql-policy-pxc-0 logs


