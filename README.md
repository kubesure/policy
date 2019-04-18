# policy

DROP TABLE policy;

CREATE TABLE `policy` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `quote_no` int(16) DEFAULT NULL,
  `receipt_no` int(10) DEFAULT NULL,
  `created_date` timestamp DEFAULT now(),
  `updated_date`  timestamp DEFAULT now(),
  PRIMARY KEY (`id`)
);