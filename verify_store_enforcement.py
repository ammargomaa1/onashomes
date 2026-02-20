import requests
import json
import sys

BASE_URL = "http://localhost:8080/api"
ADMIN_URL = "http://localhost:8080/api/admin"

def get_admin_headers():
    # Login as admin to get token
    login_payload = {
        "email": "admin@onashomes.com",
        "password": "admin123"
    }
    resp = requests.post(f"{BASE_URL}/admin/login", json=login_payload)
    if resp.status_code != 200:
        print(f"Login failed: {resp.text}")
        sys.exit(1)
    token = resp.json()["data"]["access_token"]
    return {"Authorization": f"Bearer {token}", "Content-Type": "application/json"}

headers = get_admin_headers()

def create_store_front(name, currency):
    payload = {
        "name": name,
        "currency": currency,
        "slug": name.lower().replace(" ", "-"),
        "domain": name.lower().replace(" ", "-") + ".localhost",
        "default_language": "en"
    }
    resp = requests.post(f"{ADMIN_URL}/storefronts", json=payload, headers=headers)
    if resp.status_code == 201:
        return resp.json()["data"]
    
    # If 400 and "already taken", try to fetch existing (hacky for test)
    # Ideally search by slug/name. 
    # For now, let's just make the name unique.
    return None

def create_store_front_unique(name, currency):
    import time
    timestamp = int(time.time())
    unique_name = f"{name} {timestamp}"
    return create_store_front(unique_name, currency)

def create_customer(first_name, phone, store_front_id):
    payload = {
        "first_name": first_name,
        "last_name": "Test",
        "email": f"{first_name.lower()}@test.com",
        "phone": phone,
        "store_front_id": store_front_id
    }
    resp = requests.post(f"{ADMIN_URL}/customers", json=payload, headers=headers)
    if resp.status_code != 201:
        print(f"Failed to create customer {first_name}: {resp.text}")
        return None
    return resp.json()["data"]

def create_customer_unique(prefix, store_front_id):
    import time
    import random
    unique_name = f"{prefix}{int(time.time())}"
    # Generate valid phone: 010 + 8 digits
    random_digits = "".join([str(random.randint(0, 9)) for _ in range(8)])
    phone = f"010{random_digits}"
    return create_customer(unique_name, phone, store_front_id)

def test_search_filtering():
    print("--- Testing Search Filtering ---")
    
    # 1. Create two store fronts
    sf1 = create_store_front_unique("Store A", "USD")
    sf2 = create_store_front_unique("Store B", "EUR")
    
    if not sf1 or not sf2:
        print("Failed to create store fronts")
        return
    print(f"SF1 ID: {sf1['id']}, SF2 ID: {sf2['id']}")

    # 2. Create customers for each
    c1 = create_customer_unique("CustA", sf1["id"])
    c2 = create_customer_unique("CustB", sf2["id"])

    if not c1 or not c2:
        print("Failed to create customers")
        return
    print(f"C1 ID: {c1['id']} (SF={c1['store_front_id']}), C2 ID: {c2['id']} (SF={c2['store_front_id']})")

    # 3. Search without filter (Should find both if query matches both, or specific if query matches one)
    # Search for "Cust"
    resp = requests.get(f"{ADMIN_URL}/customers/search?q=Cust", headers=headers)
    print(f"Search 'Cust': found {len(resp.json()['data'])} (Expected >= 2)")
    
    # 4. Search with SF1 filter
    resp = requests.get(f"{ADMIN_URL}/customers/search?q=Cust&store_front_id={sf1['id']}", headers=headers)
    data = resp.json()['data']
    print(f"Search 'Cust' + SF1: found {len(data)}")
    
    found_c1 = any(c['id'] == c1['id'] for c in data)
    found_c2 = any(c['id'] == c2['id'] for c in data)
    
    if found_c1 and not found_c2:
        print("SUCCESS: Filtered properly (Found CustA, Filtered out CustB)")
    else:
        print(f"FAILURE: Filter failed. Found C1: {found_c1}, Found C2: {found_c2}")

if __name__ == "__main__":
    test_search_filtering()
