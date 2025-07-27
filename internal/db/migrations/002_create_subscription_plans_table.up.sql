CREATE TABLE subscription_plans (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    plan_name VARCHAR(255) NOT NULL,
    duration INTEGER NOT NULL CHECK (duration > 0), -- number of days
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_subscription_plans_product_id ON subscription_plans(product_id);
CREATE INDEX idx_subscription_plans_created_at ON subscription_plans(created_at);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_subscription_plans_updated_at BEFORE UPDATE
    ON subscription_plans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
