CREATE TABLE `products`
(
    id       VARCHAR(36) NOT NULL,
    PRIMARY KEY (id)

) CHARACTER SET utf8mb4
  COLLATE utf8mb4_bin;

CREATE TABLE `stock_products`
(
    id       VARCHAR(36) NOT NULL,
    product_id       VARCHAR(36) NOT NULL,
    quantity INT(11) NOT NULL,

    PRIMARY KEY (id)

) CHARACTER SET utf8mb4
  COLLATE utf8mb4_bin;