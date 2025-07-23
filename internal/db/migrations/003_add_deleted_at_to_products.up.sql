-- Add deleted_at column for soft deletes
ALTER TABLE products ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

-- Create index for deleted_at
CREATE INDEX idx_products_deleted_at ON products(deleted_at);
