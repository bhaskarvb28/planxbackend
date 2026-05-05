CREATE TABLE vendor_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,

    shop_name TEXT NOT NULL,
    gstin TEXT NOT NULL,

    status TEXT DEFAULT 'pending', -- pending, approved, rejected

    created_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE vendor_applications ENABLE ROW LEVEL SECURITY;

-- INSERT
CREATE POLICY "Only normal users can apply as vendor"
ON vendor_applications
FOR INSERT
WITH CHECK (
    auth.uid() = user_id
    AND EXISTS (
        SELECT 1
        FROM profiles p
        JOIN roles r ON p.role_id = r.id
        WHERE p.id = auth.uid()
        AND r.name = 'user'
    )
);

-- SELECT
CREATE POLICY "User can view own applications"
ON vendor_applications
FOR SELECT
USING (auth.uid() = user_id);

-- Only ONE "pending" application per user
CREATE UNIQUE INDEX vendor_applications_one_pending
ON vendor_applications (user_id)
WHERE status = 'pending';