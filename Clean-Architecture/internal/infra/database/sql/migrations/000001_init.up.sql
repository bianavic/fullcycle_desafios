CREATE TABLE orders (
                        id VARCHAR(255) NOT NULL,
                        price FLOAT NOT NULL,
                        tax FLOAT NOT NULL,
                        final_price FLOAT NOT NULL,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        PRIMARY KEY (id)
);