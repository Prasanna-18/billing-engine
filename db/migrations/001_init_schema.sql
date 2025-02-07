CREATE TABLE loans (
    id SERIAL PRIMARY KEY,
    amount DECIMAL(10,2) NOT NULL,
    interest_rate DECIMAL(5,2) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER REFERENCES loans(id),
    week_number INTEGER NOT NULL,
    due_amount DECIMAL(10,2) NOT NULL,
    due_date TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(loan_id, week_number)
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER REFERENCES loans(id),
    schedule_id INTEGER REFERENCES schedules(id),
    amount DECIMAL(10,2) NOT NULL,
    payment_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_loan_id ON schedules(loan_id);
CREATE INDEX idx_loan_payments ON payments(loan_id);
CREATE INDEX idx_schedule_payments ON payments(schedule_id);