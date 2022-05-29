BEGIN;
	CREATE TABLE public.orders (
		order_uid varchar NOT NULL,
		track_number varchar NOT NULL,
		entry varchar NOT NULL,
		locale varchar NOT NULL,
		internal_signature varchar NULL,
		customer_id varchar NOT NULL,
		delivery_service varchar NOT NULL,
		shardkey varchar NULL,
		sm_id int4 NOT NULL,
		date_created timestamp NOT NULL,
		oof_shard varchar NOT NULL,
		id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
		delivery jsonb NOT NULL,
		CONSTRAINT orders_pk PRIMARY KEY (order_uid),
		CONSTRAINT orders_un UNIQUE (track_number)
	);
	CREATE INDEX orders_order_uid_idx ON public.orders USING btree (order_uid);

	CREATE TABLE public.items (
		chrt_id int4 NULL,
		track_number varchar NULL,
		price int4 NULL,
		rid varchar NULL,
		"name" varchar NULL,
		sale int4 NULL,
		"size" varchar NULL,
		total_price int4 NULL,
		nm_id int4 NULL,
		brand varchar NULL,
		status int4 NULL
	);
	CREATE INDEX items_track_number_idx ON public.items USING btree (track_number);

	ALTER TABLE public.items ADD CONSTRAINT items_fk FOREIGN KEY (track_number) REFERENCES public.orders(track_number);

	CREATE TABLE public.payment (
		"transaction" varchar NULL,
		request_id varchar NULL,
		currency varchar NULL,
		provider varchar NULL,
		amount int4 NULL,
		payment_dt int4 NULL,
		bank varchar NULL,
		delivery_cost int4 NULL,
		goods_total int4 NULL,
		custom_fee int4 NULL,
		id int4 NOT NULL GENERATED ALWAYS AS IDENTITY,
		CONSTRAINT payment_un UNIQUE (transaction)
	);
	CREATE INDEX payment_transaction_idx ON public.payment USING btree (transaction);

	ALTER TABLE public.payment ADD CONSTRAINT payment_fk FOREIGN KEY ("transaction") REFERENCES public.orders(order_uid);

COMMIT;