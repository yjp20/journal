CREATE TABLE IF NOT EXISTS media (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	description TEXT NOT NULL,
	media_type TEXT NOT NULL,
	cart BOOLEAN NOT NULL,
	completed BOOLEAN NOT NULL,
	progress FLOAT NOT NULL,

	rating FLOAT,
	notes TEXT,
	related_link TEXT,
	comments TEXT,
	completed_date TIMESTAMP(0),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS todos (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	description TEXT NOT NULL,
	private BOOLEAN NOT NULL DEFAULT false,
	cart BOOLEAN NOT NULL,
	completed BOOLEAN NOT NULL,
	progress FLOAT NOT NULL,
	blocked BOOLEAN NOT NULL,

	recur INT,
	due_date TIMESTAMP(0),
	completed_date TIMESTAMP(0),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS feed_sources (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	description TEXT NOT NULL,
	url TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS feed_items (
	id BIGSERIAL PRIMARY KEY NOT NULL,
	description TEXT NOT NULL,
	related_link TEXT NOT NULL UNIQUE,
	comments TEXT,
	media_type TEXT NOT NULL,
	added BOOLEAN NOT NULL,
	source_id BIGSERIAL,
	post_date TIMESTAMP(0),

	CONSTRAINT fk_source
		FOREIGN KEY(source_id)
			REFERENCES feed_sources(id)
			ON DELETE SET NULL
);
