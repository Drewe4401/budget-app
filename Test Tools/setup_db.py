#!/usr/bin/env python
import psycopg2
import sys

# Connection URI
URI = "postgres://admin:admin@localhost:5432/budgetdb?sslmode=disable"

def create_tables(cursor):
    create_users_table = """
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(255) UNIQUE NOT NULL, 
        password VARCHAR(255) NOT NULL,         
        permissions VARCHAR(50) NOT NULL DEFAULT 'user'
    );
    """
    create_budgets_table = """
    CREATE TABLE IF NOT EXISTS budgets (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        amount NUMERIC(10,2) NOT NULL,
        description TEXT,
        period VARCHAR(20),
        user_id INTEGER NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    """
    create_charges_table = """
    CREATE TABLE IF NOT EXISTS charges (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        amount NUMERIC(10,2) NOT NULL,
        charge_type VARCHAR(50) NOT NULL,
        periodical VARCHAR(20),
        user_id INTEGER NOT NULL,
        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );
    """
    create_shares_table = """
    CREATE TABLE IF NOT EXISTS shares (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        user_share_id INTEGER NOT NULL,
        access TEXT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (user_share_id) REFERENCES users(id) ON DELETE CASCADE
    );
    """
    
    queries = [
        create_users_table,
        create_budgets_table,
        create_charges_table,
        create_shares_table,
    ]
    
    for query in queries:
        cursor.execute(query)

def insert_filler_data(cursor):
    # Insert sample users.
    # In a real application, hash passwords securely.
    filler_users = """
    INSERT INTO users (username, password, permissions)
    VALUES
      ('admin', 'admin', 'admin'),
      ('alice', 'alice', 'user'),
      ('bob', 'bob', 'user')
    ON CONFLICT (username) DO NOTHING;
    """
    
    # Insert sample budgets using subqueries to fetch the user_id by username.
    filler_budgets = """
    INSERT INTO budgets (name, amount, category, period, user_id)
    VALUES
      ('Groceries', 500.00, 'grocery', 'monthly', (SELECT id FROM users WHERE username = 'alice')),
      ('Rent', 1200.00, 'rent', 'monthly', (SELECT id FROM users WHERE username = 'alice')),
      ('Vacation Fund', 2000.00, 'vacation', 'yearly', (SELECT id FROM users WHERE username = 'bob'))
    ON CONFLICT DO NOTHING;
    """
    
    # Insert sample charges using subqueries.
    filler_charges = """
    INSERT INTO charges (name, amount, category, periodical, user_id)
    VALUES
      ('Electricity Bill', 75.00, 'Utility', 'monthly', (SELECT id FROM users WHERE username = 'alice')),
      ('Water Bill', 30.00, 'Utility', 'monthly', (SELECT id FROM users WHERE username = 'alice')),
      ('Gym Membership', 40.00, 'Subscription', 'monthly', (SELECT id FROM users WHERE username = 'bob'))
    ON CONFLICT DO NOTHING;
    """
    
    # Insert sample shares using subqueries.
    filler_shares = """
    INSERT INTO shares (user_id, user_share_id, access)
    VALUES
      ((SELECT id FROM users WHERE username = 'alice'), (SELECT id FROM users WHERE username = 'bob'), 'read-only'),
      ((SELECT id FROM users WHERE username = 'bob'), (SELECT id FROM users WHERE username = 'alice'), 'edit')
    ON CONFLICT DO NOTHING;
    """
    
    filler_queries = [filler_budgets, filler_charges, filler_shares]
    
    for query in filler_queries:
        cursor.execute(query)

def main():
    try:
        conn = psycopg2.connect(URI)
    except Exception as e:
        print("Error connecting to the database:", e)
        sys.exit(1)
    
    try:
        cursor = conn.cursor()
        print("Connected to the database.")
        
        # Create tables
        create_tables(cursor)
        conn.commit()
        print("Tables created successfully.")
        
        # Insert filler data
        insert_filler_data(cursor)
        conn.commit()
        print("Filler data inserted successfully.")
        
    except Exception as e:
        conn.rollback()
        print("An error occurred:", e)
    finally:
        cursor.close()
        conn.close()

if __name__ == "__main__":
    main()
