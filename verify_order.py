
import requests
import sys

BASE_URL = "http://localhost:8080/api"

def main():
    # 1. Login
    print("Logging in...")
    resp = requests.post(f"{BASE_URL}/admin/login", json={
        "email": "admin@onashomes.com",
        "password": "admin123"
    })
    if resp.status_code != 200:
        print(f"Login failed: {resp.status_code} {resp.text}")
        sys.exit(1)
    
    token = resp.json()["data"]["access_token"]
    headers = {"Authorization": f"Bearer {token}"}
    print("Login successful.")

    # 2. List Orders with Customer Filter
    # Assuming we have orders from previous steps or seeding.
    # We'll list all orders first to find a customer ID
    print("Fetching recent orders...")
    resp = requests.get(f"{BASE_URL}/admin/orders?limit=5", headers=headers)
    if resp.status_code != 200:
        print(f"List orders failed: {resp.text}")
        sys.exit(1)
        
    orders = resp.json()["data"]
    if not orders:
        print("No orders found. Cannot verify auto-fill API.")
        # Optional: Create one
        return

    customer_id = orders[0].get("customer_id")
    if not customer_id:
        print("First order has no customer_id. Skipping.")
        return
        
    print(f"Testing filter for Customer ID: {customer_id}")
    
    # 3. Filter by Customer ID
    resp = requests.get(f"{BASE_URL}/admin/orders?customer_id={customer_id}&limit=1&sort=created_at&order=desc", headers=headers)
    if resp.status_code != 200:
        print(f"Filter by customer failed: {resp.text}")
        sys.exit(1)
        
    filtered_orders = resp.json()["data"]
    if not filtered_orders:
        print("Filtered list is empty!")
        sys.exit(1)
        
    fetched_order = filtered_orders[0]
    if fetched_order["customer_id"] != customer_id:
        print(f"Mismatch! Expected customer_id {customer_id}, got {fetched_order['customer_id']}")
        sys.exit(1)
        
    print("SUCCESS: Auto-fill API verified (Order filtering by Customer ID works).")
    if fetched_order.get("address"):
        print(" - Address data is present.")
    else:
        print(" - Address data is MISSING (might be expected if old order).")

if __name__ == "__main__":
    main()
