package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=onas_ecommerce port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 1. Check Roles
	var roles []map[string]interface{}
	db.Table("roles").Find(&roles)
	fmt.Println("Roles:", roles)

	// 2. Check Permission orders.edit
	var perm map[string]interface{}
	db.Table("permissions").Where("name = ?", "orders.edit").First(&perm)
	fmt.Println("Permission orders.edit:", perm)

	if perm == nil {
		fmt.Println("Creating orders.edit permission...")
		perm = map[string]interface{}{
			"name":        "orders.edit",
			"description": "Edit orders",
		}
		db.Table("permissions").Create(&perm)
	}

	// 3. Check Role Permission
	var rolePerm map[string]interface{}
	db.Table("role_permissions").Where("role_id = ? AND permission_id = ?", 1, perm["id"]).First(&rolePerm)
	fmt.Println("Role Permission (Role 1 -> orders.edit):", rolePerm)

	if rolePerm == nil {
		fmt.Println("Assigning orders.edit to Role 1...")
		db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", 1, perm["id"])
	} else {
		fmt.Println("Role 1 already has orders.edit")
	}

	// 4. Check SuperAdmin user role
	var admin map[string]interface{}
	db.Table("admins").Where("email = ?", "superadmin@onas.com").First(&admin)
	fmt.Println("SuperAdmin User:", admin)
}
