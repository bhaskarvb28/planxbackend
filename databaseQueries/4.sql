CREATE TABLE engineers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT,
    specialization TEXT,
    city TEXT,

    created_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE engineers ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Only admin can insert engineers"
ON engineers
FOR INSERT
WITH CHECK (
    EXISTS (
        SELECT 1
        FROM profiles p
        JOIN roles r ON p.role_id = r.id
        WHERE p.id = auth.uid()
        AND r.name = 'admin'
    )
);

CREATE POLICY "Users and Admin can view engineers"
ON engineers
FOR SELECT
USING (
    EXISTS (
        SELECT 1
        FROM profiles p
        JOIN roles r ON p.role_id = r.id
        WHERE p.id = auth.uid()
        AND r.name IN ('user', 'admin')
    )
);

INSERT INTO engineers (name, phone, email, specialization, city)
VALUES
('Rahul Sharma', '9876543210', 'rahul@example.com', 'Civil Engineer', 'Bangalore'),
('Anita Verma', '9123456780', 'anita@example.com', 'Interior Engineer', 'Mumbai'),
('Kiran Patel', '9988776655', 'kiran@example.com', 'Structural Engineer', 'Ahmedabad'),
('Sneha Reddy', '9012345678', 'sneha@example.com', 'Interior Designer', 'Hyderabad'),
('Arjun Mehta', '9090909090', 'arjun@example.com', 'Civil Engineer', 'Delhi');