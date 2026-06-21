CREATE TABLE enrollments (
	id UUID PRIMARY KEY,

	student_id UUID NOT NULL,
	course_id UUID NOT NULL,

	created_at TIMESTAMP NOT NULL,

	FOREIGN KEY(student_id)
	REFERENCES students(id),

	FOREIGN KEY(course_id)
	REFERENCES courses(id),

	UNIQUE(student_id, course_id)
);
