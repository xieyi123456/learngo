package service

import (
	"testing"
)

// TestCompareJSON_IdenticalJSON 测试比较两个完全相同的JSON
func TestCompareJSON_IdenticalJSON(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{"name":"Alice","age":30,"address":{"city":"New York"}}`
	json2 := `{"name":"Alice","age":30,"address":{"city":"New York"}}`

	diff, err := service.CompareJSON(json1, json2)
	if err != nil {
		t.Errorf("比较相同JSON时出错: %v", err)
		return
	}

	if len(diff.Added) > 0 {
		t.Errorf("预期没有新增字段，但实际有: %v", diff.Added)
	}

	if len(diff.Removed) > 0 {
		t.Errorf("预期没有移除字段，但实际有: %v", diff.Removed)
	}

	if len(diff.Changed) > 0 {
		t.Errorf("预期没有变更字段，但实际有: %v", diff.Changed)
	}
}

// TestCompareJSON_DifferentValues 测试比较有值差异的两个JSON
func TestCompareJSON_DifferentValues(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{"name":"Alice","age":30,"active":true}`
	json2 := `{"name":"Alice","age":31,"active":false}`

	diff, err := service.CompareJSON(json1, json2)
	if err != nil {
		t.Errorf("比较有值差异的JSON时出错: %v", err)
		return
	}

	if len(diff.Added) > 0 {
		t.Errorf("预期没有新增字段，但实际有: %v", diff.Added)
	}

	if len(diff.Removed) > 0 {
		t.Errorf("预期没有移除字段，但实际有: %v", diff.Removed)
	}

	if len(diff.Changed) != 2 {
		t.Errorf("预期有2个变更字段，但实际有%d个: %v", len(diff.Changed), diff.Changed)
	}

	if _, exists := diff.Changed["age"]; !exists {
		t.Errorf("预期'age'字段发生变更，但实际没有")
	}

	if _, exists := diff.Changed["active"]; !exists {
		t.Errorf("预期'active'字段发生变更，但实际没有")
	}
}

// TestCompareJSON_AddRemoveFields 测试比较有新增和移除字段的两个JSON
func TestCompareJSON_AddRemoveFields(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{"name":"Alice","age":30,"department":"Engineering"}`
	json2 := `{"name":"Alice","email":"alice@example.com","location":"New York"}`

	diff, err := service.CompareJSON(json1, json2)
	if err != nil {
		t.Errorf("比较有新增和移除字段的JSON时出错: %v", err)
		return
	}

	// 检查新增字段
	expectedAdded := []string{"email", "location"}
	if len(diff.Added) != len(expectedAdded) {
		t.Errorf("预期有%d个新增字段，但实际有%d个: %v", len(expectedAdded), len(diff.Added), diff.Added)
	}

	for _, field := range expectedAdded {
		found := false
		for _, added := range diff.Added {
			if added == field {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("预期'%s'字段被新增，但实际没有", field)
		}
	}

	// 检查移除字段
	expectedRemoved := []string{"age", "department"}
	if len(diff.Removed) != len(expectedRemoved) {
		t.Errorf("预期有%d个移除字段，但实际有%d个: %v", len(expectedRemoved), len(diff.Removed), diff.Removed)
	}

	for _, field := range expectedRemoved {
		found := false
		for _, removed := range diff.Removed {
			if removed == field {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("预期'%s'字段被移除，但实际没有", field)
		}
	}
}

// TestCompareJSON_NestedObjects 测试比较嵌套对象的两个JSON
func TestCompareJSON_NestedObjects(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{"person":{"name":"Alice","address":{"city":"New York","zip":"10001"}}}`
	json2 := `{"person":{"name":"Alice","address":{"city":"Boston","country":"USA"}}}`

	diff, err := service.CompareJSON(json1, json2)
	if err != nil {
		t.Errorf("比较嵌套对象的JSON时出错: %v", err)
		return
	}

	// 检查变更字段
	if _, exists := diff.Changed["person.address.city"]; !exists {
		t.Errorf("预期'person.address.city'字段发生变更，但实际没有")
	}

	// 检查新增字段
	found := false
	for _, added := range diff.Added {
		if added == "person.address.country" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("预期'person.address.country'字段被新增，但实际没有")
	}

	// 检查移除字段
	found = false
	for _, removed := range diff.Removed {
		if removed == "person.address.zip" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("预期'person.address.zip'字段被移除，但实际没有")
	}
}

// TestCompareJSON_Arrays 测试比较包含数组的两个JSON
func TestCompareJSON_Arrays(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{"name":"Alice","hobbies":["reading","swimming","coding"]}`
	json2 := `{"name":"Alice","hobbies":["reading","cycling"]}`

	diff, err := service.CompareJSON(json1, json2)
	if err != nil {
		t.Errorf("比较包含数组的JSON时出错: %v", err)
		return
	}

	// 检查数组元素变更
	if _, exists := diff.Changed["hobbies[1]"]; !exists {
		t.Errorf("预期'hobbies[1]'字段发生变更，但实际没有")
	}

	// 检查数组长度变更
	if _, exists := diff.Changed["hobbies"]; !exists {
		t.Errorf("预期'hobbies'数组长度发生变更，但实际没有")
	}
}

// TestCompareJSON_InvalidJSON 测试比较无效的JSON
func TestCompareJSON_InvalidJSON(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{"name":"Alice","age":30` // 无效的JSON（缺少右括号）
	json2 := `{"name":"Bob","age":25}`

	_, err := service.CompareJSON(json1, json2)
	if err == nil {
		t.Errorf("预期比较无效JSON时出错，但实际没有")
	}
}

// TestCompareJSON_ComplexNestedObjectsArrays 测试比较复杂的多层嵌套对象数组JSON
func TestCompareJSON_ComplexNestedObjectsArrays(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{
		"company": "TechCorp",
		"employees": [
			{
				"id": 1,
				"name": "Alice",
				"department": {
					"name": "Engineering",
					"manager": "John",
					"teamSize": 15
				},
				"skills": ["Go", "Python", "JavaScript"],
				"contact": {
					"email": "alice@techcorp.com",
					"phone": {
						"work": "123-456-7890",
						"personal": null
					}
				}
			},
			{
				"id": 2,
				"name": "Bob",
				"department": {
					"name": "Marketing",
					"manager": "Sarah",
					"teamSize": 8
				},
				"skills": ["Marketing", "Design"],
				"projects": [
					{"name": "Product Launch", "status": "active"},
					{"name": "Brand Campaign", "status": "planning"}
				]
			}
		],
		"offices": [
			{
				"location": "New York",
				"address": {
					"street": "123 Main St",
					"zip": "10001"
				}
			},
			{
				"location": "London",
				"address": {
					"street": "456 Oxford St",
					"zip": "W1D 1BS"
				}
			}
		],
		"metadata": {
			"founded": 2010,
			"industry": "Technology",
			"tags": ["software", "cloud", "AI"]
		}
	}`

	json2 := `{
		"company": "TechCorp",
		"employees": [
			{
				"id": 1,
				"name": "Alice Chen",
				"department": {
					"name": "Engineering",
					"manager": "John Doe",
					"teamSize": 20,
					"budget": "$2M"
				},
				"skills": ["Go", "TypeScript", "Docker"],
				"contact": {
					"email": "alice.chen@techcorp.com",
					"phone": {
						"work": "123-456-7890",
						"personal": "987-654-3210"
					},
					"address": {
						"city": "Boston",
						"state": "MA"
					}
				},
				"fullTime": true
			},
			{
				"id": 3,
				"name": "Charlie",
				"department": {
					"name": "Sales",
					"manager": "Mike"
				},
				"skills": ["Sales", "Negotiation"],
				"projects": [
					{"name": "Client Acquisition", "status": "active"}
				]
			}
		],
		"offices": [
			{
				"location": "New York",
				"address": {
					"street": "123 Main St",
					"zip": "10001",
					"country": "USA"
				},
				"employees": 150
			},
			{
				"location": "Berlin",
				"address": {
					"street": "789 Unter den Linden",
					"zip": "10117"
				}
			}
		],
		"metadata": {
			"founded": 2010,
			"industry": "Software",
			"tags": ["software", "cloud", "AI", "ML"],
			"rating": 4.8
		}
	}`

	diff, err := service.CompareJSON(json1, json2)
	if err != nil {
		t.Errorf("比较复杂嵌套JSON时出错: %v", err)
		return
	}

	// 调试信息：打印实际的diff结果
	t.Logf("实际的变更字段: %v", diff.Changed)
	t.Logf("实际的新增字段: %v", diff.Added)
	t.Logf("实际的移除字段: %v", diff.Removed)

	// 检查预期的变更
	expectedChanges := []string{
		"employees[0].name",
		"employees[0].department.manager",
		"employees[0].department.teamSize",
		"employees[0].skills[1]",
		"employees[0].skills[2]",
		"employees[0].contact.email",
		"employees[0].contact.phone.personal",
		"employees[1].department.manager",
		"employees[1].department.name",
		"employees[1].id",
		"employees[1].name",
		"employees[1].projects",
		"employees[1].projects[0].name",
		"employees[1].skills[0]",
		"employees[1].skills[1]",
		"metadata.industry",
		"metadata.tags",
		"offices[1].address.street",
		"offices[1].address.zip",
		"offices[1].location",
	}

	for _, path := range expectedChanges {
		if _, exists := diff.Changed[path]; !exists {
			t.Errorf("预期'%s'字段发生变更，但实际没有", path)
		}
	}

	// 检查预期的新增
	expectedAdded := []string{
		"employees[0].contact.address",
		"employees[0].department.budget",
		"employees[0].fullTime",
		"metadata.rating",
		"metadata.tags[3]",
		"offices[0].address.country",
		"offices[0].employees",
	}

	for _, path := range expectedAdded {
		found := false
		for _, added := range diff.Added {
			if added == path {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("预期'%s'字段被新增，但实际没有", path)
		}
	}

	// 检查预期的移除
	expectedRemoved := []string{
		"employees[1].department.teamSize",
		"employees[1].projects[1]",
	}

	for _, path := range expectedRemoved {
		found := false
		for _, removed := range diff.Removed {
			if removed == path {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("预期'%s'字段被移除，但实际没有", path)
		}
	}
}

// TestCompareJSONWithIgnore_IgnorePaths 测试使用忽略路径进行JSON比较
func TestCompareJSONWithIgnore_IgnorePaths(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{
		"name": "Alice",
		"age": 30,
		"email": "alice@example.com",
		"address": {
			"city": "New York",
			"zip": "10001",
			"country": "USA"
		},
		"hobbies": ["reading", "swimming"]
	}`

	json2 := `{
		"name": "Alice Chen",
		"age": 31,
		"email": "alice@example.com",
		"address": {
			"city": "Boston",
			"zip": "02101",
			"country": "USA"
		},
		"hobbies": ["reading", "cycling", "cooking"],
		"department": "Engineering"
	}`

	// 忽略 age, address.zip, hobbies[1] 和 department 路径
	ignorePaths := []string{"age", "address.zip", "hobbies[1]", "department"}

	diff, err := service.CompareJSONWithIgnore(json1, json2, ignorePaths)
	if err != nil {
		t.Errorf("使用忽略路径比较JSON时出错: %v", err)
		return
	}

	// 检查被忽略的字段不应该出现在结果中
	if _, exists := diff.Changed["age"]; exists {
		t.Errorf("'age'字段应该被忽略，但出现在变更中")
	}

	if _, exists := diff.Changed["address.zip"]; exists {
		t.Errorf("'address.zip'字段应该被忽略，但出现在变更中")
	}

	if _, exists := diff.Changed["hobbies[1]"]; exists {
		t.Errorf("'hobbies[1]'字段应该被忽略，但出现在变更中")
	}

	// 检查department不应该出现在新增字段中
	for _, added := range diff.Added {
		if added == "department" {
			t.Errorf("'department'字段应该被忽略，但出现在新增中")
		}
	}

	// 检查未被忽略的字段应该正常检测到
	if _, exists := diff.Changed["name"]; !exists {
		t.Errorf("'name'字段应该被检测到变更，但实际没有")
	}

	if _, exists := diff.Changed["address.city"]; !exists {
		t.Errorf("'address.city'字段应该被检测到变更，但实际没有")
	}

	// 检查数组新增的元素（hobbies[2]）应该被检测到
	found := false
	for _, added := range diff.Added {
		if added == "hobbies[2]" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("'hobbies[2]'字段应该被检测到新增，但实际没有")
	}
}

// TestCompareJSONWithIgnore_EmptyIgnorePaths 测试空忽略路径列表
func TestCompareJSONWithIgnore_EmptyIgnorePaths(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{"name":"Alice","age":30}`
	json2 := `{"name":"Bob","age":31}`

	// 空忽略路径列表
	diff1, err := service.CompareJSONWithIgnore(json1, json2, []string{})
	if err != nil {
		t.Errorf("使用空忽略路径列表比较JSON时出错: %v", err)
		return
	}

	// nil忽略路径
	diff2, err := service.CompareJSONWithIgnore(json1, json2, nil)
	if err != nil {
		t.Errorf("使用nil忽略路径比较JSON时出错: %v", err)
		return
	}

	// 应该与普通CompareJSON结果相同
	diff3, err := service.CompareJSON(json1, json2)
	if err != nil {
		t.Errorf("普通比较JSON时出错: %v", err)
		return
	}

	// 比较结果应该相同
	if len(diff1.Changed) != len(diff2.Changed) || len(diff1.Changed) != len(diff3.Changed) {
		t.Errorf("空忽略路径、nil忽略路径和普通比较的结果应该相同")
	}
}

// TestCompareJSONWithIgnore_ArrayWildcards 测试数组通配符功能
func TestCompareJSONWithIgnore_ArrayWildcards(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{
		"users": [
			{
				"id": 1,
				"name": "Alice",
				"metadata": {
					"created": "2023-01-01",
					"updated": "2023-06-01"
				}
			},
			{
				"id": 2,
				"name": "Bob",
				"metadata": {
					"created": "2023-02-01",
					"updated": "2023-07-01"
				}
			}
		],
		"products": [
			{"price": 10.99, "name": "Product A"},
			{"price": 20.99, "name": "Product B"}
		]
	}`

	json2 := `{
		"users": [
			{
				"id": 1,
				"name": "Alice Chen",
				"metadata": {
					"created": "2023-01-01",
					"updated": "2023-08-01"
				}
			},
			{
				"id": 2,
				"name": "Bob Smith",
				"metadata": {
					"created": "2023-02-01",
					"updated": "2023-09-01"
				}
			}
		],
		"products": [
			{"price": 12.99, "name": "Product A"},
			{"price": 22.99, "name": "Product B"}
		]
	}`

	// 使用数组通配符忽略所有用户的updated时间和所有产品的价格
	ignorePaths := []string{
		"users[*].metadata.updated", // 忽略所有用户的updated字段
		"products[*].price",         // 忽略所有产品的价格
	}

	diff, err := service.CompareJSONWithIgnore(json1, json2, ignorePaths)
	if err != nil {
		t.Errorf("使用数组通配符比较JSON时出错: %v", err)
		return
	}

	// 检查被通配符忽略的字段不应该出现在结果中
	ignoredFields := []string{
		"users[0].metadata.updated",
		"users[1].metadata.updated",
		"products[0].price",
		"products[1].price",
	}

	for _, field := range ignoredFields {
		if _, exists := diff.Changed[field]; exists {
			t.Errorf("'%s'字段应该被通配符忽略，但出现在变更中", field)
		}
	}

	// 检查未被忽略的字段应该正常检测到
	if _, exists := diff.Changed["users[0].name"]; !exists {
		t.Errorf("'users[0].name'字段应该被检测到变更，但实际没有")
	}

	if _, exists := diff.Changed["users[1].name"]; !exists {
		t.Errorf("'users[1].name'字段应该被检测到变更，但实际没有")
	}
}

// TestCompareJSONWithIgnore_ComplexWildcards 测试复杂的数组通配符场景
func TestCompareJSONWithIgnore_ComplexWildcards(t *testing.T) {
	service := NewJSONDiffService()
	json1 := `{
		"departments": [
			{
				"name": "Engineering",
				"employees": [
					{"name": "Alice", "salary": 80000, "bonus": 5000},
					{"name": "Bob", "salary": 75000, "bonus": 4000}
				]
			},
			{
				"name": "Marketing", 
				"employees": [
					{"name": "Charlie", "salary": 60000, "bonus": 3000}
				]
			}
		]
	}`

	json2 := `{
		"departments": [
			{
				"name": "Engineering",
				"employees": [
					{"name": "Alice Chen", "salary": 85000, "bonus": 6000},
					{"name": "Bob Smith", "salary": 78000, "bonus": 4500}
				]
			},
			{
				"name": "Marketing",
				"employees": [
					{"name": "Charlie Brown", "salary": 62000, "bonus": 3500}
				]
			}
		]
	}`

	// 忽略所有员工的薪水和奖金
	ignorePaths := []string{
		"departments[*].employees[*].salary",
		"departments[*].employees[*].bonus",
	}

	diff, err := service.CompareJSONWithIgnore(json1, json2, ignorePaths)
	if err != nil {
		t.Errorf("使用复杂数组通配符比较JSON时出错: %v", err)
		return
	}

	// 检查所有被忽略的薪水和奖金字段
	ignoredFields := []string{
		"departments[0].employees[0].salary",
		"departments[0].employees[0].bonus",
		"departments[0].employees[1].salary",
		"departments[0].employees[1].bonus",
		"departments[1].employees[0].salary",
		"departments[1].employees[0].bonus",
	}

	for _, field := range ignoredFields {
		if _, exists := diff.Changed[field]; exists {
			t.Errorf("'%s'字段应该被通配符忽略，但出现在变更中", field)
		}
	}

	// 检查员工姓名变更应该被检测到
	expectedNameChanges := []string{
		"departments[0].employees[0].name",
		"departments[0].employees[1].name",
		"departments[1].employees[0].name",
	}

	for _, field := range expectedNameChanges {
		if _, exists := diff.Changed[field]; !exists {
			t.Errorf("'%s'字段应该被检测到变更，但实际没有", field)
		}
	}
}

// TestMatchesPath_WildcardPatterns 测试路径匹配函数的各种通配符模式
func TestMatchesPath_WildcardPatterns(t *testing.T) {
	testCases := []struct {
		path     string
		pattern  string
		expected bool
		desc     string
	}{
		// 精确匹配
		{"users[0].name", "users[0].name", true, "精确匹配"},
		{"users[0].name", "users[1].name", false, "精确不匹配"},

		// 数组通配符
		{"users[0].name", "users[*].name", true, "单个数组通配符匹配"},
		{"users[1].name", "users[*].name", true, "单个数组通配符匹配不同索引"},
		{"users[0].age", "users[*].name", false, "通配符字段不匹配"},
		{"products[0].name", "users[*].name", false, "通配符路径不匹配"},

		// 多层数组通配符
		{"users[0].contacts[1].email", "users[*].contacts[*].email", true, "多层数组通配符匹配"},
		{"users[2].contacts[0].phone", "users[*].contacts[*].email", false, "多层通配符字段不匹配"},

		// 嵌套对象中的数组通配符
		{"company.employees[0].name", "company.employees[*].name", true, "嵌套数组通配符"},
		{"company.employees[5].salary", "company.employees[*].salary", true, "嵌套数组通配符大索引"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := matchesPath(tc.path, tc.pattern)
			if result != tc.expected {
				t.Errorf("matchesPath(%q, %q) = %v, 预期 %v", tc.path, tc.pattern, result, tc.expected)
			}
		})
	}
}
