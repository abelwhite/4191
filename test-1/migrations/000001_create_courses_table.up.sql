CREATE TABLE courses(
	id bigserial PRIMARY KEY,
    course_code text NOT NULL,
	course_title text NOT NULL,
	course_credit text NOT NULL,
	version integer NOT NULL DEFAULT 1,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()


			
);