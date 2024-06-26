CREATE TYPE token as (
    id int, 
    name varchar(40), 
    address varchar(50), 
    wei bigint
    );

CREATE TABLE tokens OF token(
    PRIMARY KEY(id)
    );

CREATE TYPE token_pair AS (
    id int,
    token0 token, 
    token1 token,
    );

CREATE TYPE swap_protocol AS (
    id int,
    name varchar(40), 
    factory varchar(50), 
    router varchar(50)
    );

CREATE TABLE swap_protocols OF swap_protocol(
    PRIMARY KEY(id)
);

CREATE TABLE token_pairs OF token_pair(
    PRIMARY KEY(id), 
    FOREIGN KEY(token0_id) 
    REFERENCES tokens(id) 
    ON UPDATE CASCADE 
    ON DELETE SET NULL, 
    FOREIGN KEY(token1_id) 
    REFERENCES tokens(id) 
    ON UPDATE CASCADE 
    ON DELETE SET NULL
    );


CREATE TYPE pool AS (
    id int, 
    address varchar(50), 
    pair_id token_pair, 
    protocol_id swap_protocol
    );

CREATE TABLE pools OF pool (
    PRIMARY KEY(id), 
    FOREIGN KEY(pair_id) 
    REFERENCES token_pairs(id) 
    ON UPDATE CASCADE 
    ON DELETE SET NULL, 
    FOREIGN KEY(protocol_id) 
    REFERENCES swap_protocols(id) 
    ON UPDATE CASCADE 
    ON DELETE SET NULL
    );

