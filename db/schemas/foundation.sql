
-- This file contains the structure of the database

START TRANSACTION;

CREATE TABLE IF NOT EXISTS assets (
    id int NOT NULL auto_increment,
    isin varchar(64) NOT NULL,
    name varchar(256) NOT NULL,
    PRIMARY KEY(id),
    UNIQUE KEY(isin)
);

CREATE TABLE IF NOT EXISTS currencies (
    id int NOT NULL auto_increment,
    code varchar(3) NOT NULL,
    UNIQUE KEY(code),
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS prices (
    id int NOT NULL auto_increment,
    asset_id int NOT NULL,
    price decimal(12, 2) NOT NULL,
    currency_id int NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY (asset_id) REFERENCES assets(id),
    FOREIGN KEY (currency_id) REFERENCES currencies(id)
);

CREATE TABLE IF NOT EXISTS investors (
    id int NOT NULL auto_increment,
    name varchar(64) NOT NULL,
    /* plain text pass should NEVER be done in production, but done here for demo purposes */
    pass varchar(64) NOT NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS asset_holdings (
    id int NOT NULL auto_increment,
    investor_id int NOT NULL,
    asset_id int NOT NULL,
    units int NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY (investor_id) REFERENCES investors(id),
    FOREIGN KEY (asset_id) REFERENCES assets(id)
);

CREATE TABLE IF NOT EXISTS currency_holdings (
    id int NOT NULL auto_increment,
    investor_id int NOT NULL,
    currency_id int NOT NULL,
    /* amount is in the currency's smallest unit value. e.g 100 = 1 pound */
    amount int NOT NULL,
    PRIMARY KEY(id),
    FOREIGN KEY (investor_id) REFERENCES investors(id),
    FOREIGN KEY (currency_id) REFERENCES currencies(id)
);

COMMIT;