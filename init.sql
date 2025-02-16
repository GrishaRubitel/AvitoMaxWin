CREATE TABLE users (
	login varchar NOT NULL,
	pass_hash varchar NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (login),
	CONSTRAINT users_unique UNIQUE (login)
);

CREATE TABLE users_cash (
	login varchar NOT NULL,
	cash int4 DEFAULT 1000 NULL,
	CONSTRAINT users_cash_pk PRIMARY KEY (login),
	CONSTRAINT users_cash_check CHECK (cash >= 0),
	CONSTRAINT users_cash_users_fk FOREIGN KEY (login) REFERENCES users(login) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE "catalog" (
	item varchar NOT NULL,
	price int4 NULL,
	CONSTRAINT catalog_pk PRIMARY KEY (item),
	CONSTRAINT catalog_check CHECK (price >= 0)
);

CREATE TABLE inventory (
	login varchar NOT NULL,
	item varchar NOT NULL,
	quantity int4 DEFAULT 1 NOT NULL
	CONSTRAINT inventory_check CHECK (quantity > 0),
	CONSTRAINT inventory_pk PRIMARY KEY (login,item),
	CONSTRAINT inventory_users_fk FOREIGN KEY (login) REFERENCES users(login) ON DELETE SET DEFAULT ON UPDATE CASCADE,
	CONSTRAINT inventory_catalog_fk FOREIGN KEY (item) REFERENCES "catalog"(item) ON DELETE SET DEFAULT ON UPDATE CASCADE
);

CREATE TABLE transactions (
	id serial4 NOT NULL,
	sender varchar NULL,
	recipient varchar NULL,
	amount int4 NOT NULL,
	transaction_type varchar(50) NOT NULL,
	item varchar NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT transactions_amount_check CHECK ((amount > 0)),
	CONSTRAINT transactions_pkey PRIMARY KEY (id),
	CONSTRAINT transactions_item_fkey FOREIGN KEY (item) REFERENCES "catalog"(item) ON DELETE SET NULL,
	CONSTRAINT transactions_recipient_fkey FOREIGN KEY (recipient) REFERENCES users(login) ON DELETE SET NULL,
	CONSTRAINT transactions_sender_fkey FOREIGN KEY (sender) REFERENCES users(login) ON DELETE SET NULL
);

INSERT INTO "catalog" (item, price) VALUES
('t-shirt', 80),
('cup', 20),
('book', 50),
('pen', 10),
('powerbank', 200),
('hoody', 300),
('umbrella', 200),
('socks', 10),
('wallet', 50),
('pink-hoody', 500),
('cheesborg', 2147483647);

INSERT INTO users (login, pass_hash) VALUES
('joe_peach', '$2a$10$JB7iUgigSC7t.zbh2aPGFuVpI57.ILjORdYoo6qaSkMNNWuB92P6e'),
('deadp47', '$2a$10$FAF6DQ4XBU0IxuGHNhD03OCSvDRENlM3OsmkaafV6S2elX8Tk1kZq'),
('bobs', '$2a$10$goOCK.uZNu/owKLHhe.52uiZ672KA/7Fqvc2m7yMSSGM05NKSgNLm'),
('bebs', '$2a$10$BfFk5MsLB4LSNeOBwCtwjus5f0LsibD1Fo84GZZK4GzJUKDbuUSRS'),
('babs', '$2a$10$8l8IGcA4te6f7obP505lFOSblm5THp8wfYR9Z9Binf8bAPQX7MM2u');

INSERT INTO users_cash (login, cash) VALUES
('joe_peach', 1000), ('deadp47', 1000), ('bobs', 1000), ('bebs', 1000), ('babs', 1000);

INSERT INTO transactions (sender, recipient, amount, transaction_type, created_at) VALUES
('joe_peach', 'deadp47', 200, 'transfer', NOW()),
('deadp47', 'joe_peach', 250, 'transfer', NOW()),
('bobs', 'bebs', 500, 'transfer', NOW()),
('bebs', 'bobs', 150, 'transfer', NOW()),
('babs', 'joe_peach', 400, 'transfer', NOW()),
('joe_peach', 'babs', 220, 'transfer', NOW()),
('deadp47', 'bobs', 110, 'transfer', NOW()),
('bobs', 'babs', 180, 'transfer', NOW()),
('bebs', 'deadp47', 100, 'transfer', NOW()),
('babs', 'bebs', 180, 'transfer', NOW());

INSERT INTO transactions (sender, amount, transaction_type, item, created_at) VALUES
('joe_peach', 80, 'purchase', 't-shirt', NOW()),
('deadp47', 20, 'purchase', 'cup', NOW()),
('bobs', 50, 'purchase', 'book', NOW()), 
('bebs', 10, 'purchase', 'pen', NOW()),
('babs', 200, 'purchase', 'powerbank', NOW());

INSERT INTO inventory (login, item, quantity) VALUES
('joe_peach', 't-shirt', 1), ('deadp47', 'cup', 1), ('bobs', 'book', 1), ('bebs', 'pen', 413), ('babs', 'powerbank', 7);

create or replace function get_user_cash(username varchar) 
returns int as 
$$
declare
    cash int;
begin
    select uc.cash into cash 
    from users_cash uc
    where uc.login = username;

    if cash is null then
        return 0;
    end if;

    return cash;
end;
$$
language plpgsql;



create or replace function send_coins(sender varchar, recipient varchar, amount int)
returns void as 
$$
declare
    sender_balance int;
begin
    begin
        select cash into sender_balance 
		from users_cash
		where login = sender 
		for update;
    
        if sender_balance < amount then
            raise exception 'not enough money';
        end if;

        update users_cash
		set cash = cash - amount 
		where login = sender;
        
		update users_cash
		set cash = cash + amount
		where login = recipient;

        insert into transactions (sender, recipient, amount, transaction_type) 
        values (sender, recipient, amount, 'transfer');
    exception
        when others then
            raise;
    end;
end;
$$
language plpgsql;



create or replace function buy_item(par_login varchar, par_item varchar)
returns void as 
$$
declare
    item_price int;
    user_balance int;
    existing_quantity int;
begin
    select cash into user_balance
	from users_cash
	where login = par_login for update;

    select price into item_price 
	from "catalog" 
	where item = par_item;

    if item_price is null then
        raise exception 'item not in stock';
    end if;

    if user_balance < item_price then
        raise exception 'not enough money';
    end if;

    update users_cash
	set cash = cash - item_price
	where login = par_login;

    select quantity into existing_quantity
	from inventory
	where login = par_login and item = par_item;

    if existing_quantity is not null then
        update inventory set quantity = quantity + 1 where login = par_login and item = par_item;
    else
        insert into inventory (login, item, quantity) values (par_login, par_item, 1);
    end if;

    insert into transactions (sender, recipient, amount, transaction_type, item) 
    values (par_login, null, item_price, 'purchase', par_item);

	exception
		when others then
			raise;
	end;
$$
language plpgsql;
