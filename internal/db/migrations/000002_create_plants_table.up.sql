CREATE TABLE IF NOT EXISTS plants (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    acquisition_date TIMESTAMP NOT NULL,
    location VARCHAR(255),
    care_frequency INT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);