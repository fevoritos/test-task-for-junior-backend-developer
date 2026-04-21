CREATE TABLE IF NOT EXISTS recurrence_tasks (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    
    recur_type TEXT NOT NULL, 

    interval_days INT,         
    day_of_month INT,           
    specific_dates TEXT[],      
    parity TEXT,                

    start_date DATE NOT NULL,   
    end_date DATE,              
    last_generated_date DATE,   
    
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tasks (
    id BIGSERIAL PRIMARY KEY,
    recur_id BIGINT REFERENCES recurrence_tasks(id) ON DELETE SET NULL,
    
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'todo', 
    
    scheduled_at TIMESTAMPTZ NOT NULL, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_schedule_at ON tasks(scheduled_at);
