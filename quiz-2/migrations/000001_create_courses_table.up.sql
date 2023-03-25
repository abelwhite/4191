CREATE TABLE courses(
    schools_id bigserial PRIMARY KEY,
    name text NOT NULL, 
	level VARCHAR(100) NOT NULL,
	contact VARCHAR(100) NOT NULL,
	phone VARCHAR(25) NOT NULL,
	email VARCHAR(25) NOT NULL,
	website VARCHAR(100) NOT NULL,
	address VARCHAR(100) NOT NULL,
	mode VARCHAR(100) NOT NULL,
	Version serial,
	 created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);