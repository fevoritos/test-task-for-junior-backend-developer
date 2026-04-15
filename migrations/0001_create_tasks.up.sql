CREATE TABLE IF NOT EXISTS recurrence_rules (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    
    rule_type TEXT NOT NULL, 

    interval_days INT,         
    day_of_month INT,           
    specific_dates DATE[],      
    parity TEXT,                

    start_date DATE NOT NULL,   
    end_date DATE,              
    last_generated_date DATE,   
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    recur_id BIGINT REFERENCES recurrence_rules(id) ON DELETE SET NULL,
    
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'todo', 
    
    due_date DATE NOT NULL, 
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_due_date ON tasks(due_date);
