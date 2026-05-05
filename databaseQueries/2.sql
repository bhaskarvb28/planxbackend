-- Credits Table
CREATE TABLE user_credits (
    user_id UUID PRIMARY KEY REFERENCES profiles(id) ON DELETE CASCADE,
    credits INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Enable RLS
ALTER TABLE user_credits ENABLE ROW LEVEL SECURITY;

-- Policy 1 : User can view own credits
CREATE POLICY "User can view own credits"
ON user_credits
FOR SELECT
USING (auth.uid() = user_id)