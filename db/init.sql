CREATE TABLE users
(
    username text primary key,
    password text NOT NULL,
    role     text DEFAULT 'buyer',
    deposit  int  DEFAULT 0,
    CONSTRAINT chk_role CHECK (role IN ('buyer', 'seller'))
);


CREATE TABLE products
(
    name      text primary key,
    seller_id text,
    CONSTRAINT fk_seller
        FOREIGN KEY (seller_id)
            REFERENCES users (username) ON DELETE CASCADE
);

CREATE FUNCTION is_seller() RETURNS trigger AS
$isseller$
DECLARE
    user_role text;
BEGIN
    select role into user_role from users where username = NEW.seller_id;
    if user_role != 'seller' THEN
        RAISE EXCEPTION 'role is not seller';
    end if;
    RETURN NEW;
END;
$isseller$ LANGUAGE plpgsql;

CREATE TRIGGER check_role
    BEFORE INSERT OR UPDATE
    ON products
    FOR EACH ROW
EXECUTE PROCEDURE is_seller();

CREATE TABLE transactions
(
    id           serial primary key,
    product_name text,
    username     text,
    amount       INT default 1,
    price        INT,
    CONSTRAINT fk_product_name
        FOREIGN KEY (product_name)
            REFERENCES products (name),
    CONSTRAINT fk_username
        FOREIGN KEY (username)
            REFERENCES users (username) ON DELETE CASCADE
);


CREATE FUNCTION update_inventory() RETURNS trigger AS
$update_inventory$
DECLARE
    inventory_amount int;
    product_price    double precision;
    user_deposit     int;
BEGIN
    SELECT amount, price into inventory_amount,product_price from inventory where product_name = NEW.product_name;
    if NEW.amount > inventory_amount then
        RAISE EXCEPTION 'amount is larger than inventory';
    end if;
    SELECT deposit into user_deposit from users where username = NEW.username;
    NEW.price = CEILING((product_price / 5)) * 5;
    if NEW.amount * NEW.price > user_deposit THEN
        RAISE EXCEPTION 'cost is higher than deposit';
    end if;

    UPDATE inventory SET amount = amount - new.amount WHERE product_name = NEW.product_name;
    UPDATE users SET deposit = deposit - (NEW.amount * NEW.price) WHERE username = NEW.username;
    RETURN NEW;
END;
$update_inventory$ LANGUAGE plpgsql;

CREATE TRIGGER check_update
    BEFORE INSERT
    ON transactions
    FOR EACH ROW
EXECUTE PROCEDURE update_inventory();


CREATE TABLE inventory
(
    id           serial primary key,
    product_name text unique,
    amount       INT,
    price        int,
    CONSTRAINT fk_product_name
        FOREIGN KEY (product_name)
            REFERENCES products (name) on delete cascade
);
