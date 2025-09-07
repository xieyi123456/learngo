package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// JSONDiffResult 表示两个JSON之间的差异结果
type JSONDiffResult struct {
	Added   []string          `json:"added"`   // 在第二个JSON中新增的键路径
	Removed []string          `json:"removed"` // 在第一个JSON中存在但在第二个中不存在的键路径
	Changed map[string]string `json:"changed"` // 值发生变化的键路径和对应的变更信息
}

// JSONDiffService 提供JSON差异比较的服务接口
type JSONDiffService interface {
	CompareJSON(json1, json2 string) (JSONDiffResult, error)
	CompareJSONWithIgnore(json1, json2 string, ignorePaths []string) (JSONDiffResult, error)
}

// jsonDiffServiceImpl JSON差异比较服务的具体实现
type jsonDiffServiceImpl struct{}

// NewJSONDiffService 创建一个新的JSON差异比较服务实例
func NewJSONDiffService() JSONDiffService {
	return &jsonDiffServiceImpl{}
}

// CompareJSON 比较两个JSON字符串并返回它们之间的差异
func (s *jsonDiffServiceImpl) CompareJSON(json1, json2 string) (JSONDiffResult, error) {
	return s.CompareJSONWithIgnore(json1, json2, nil)
}

// CompareJSONWithIgnore 比较两个JSON字符串并返回它们之间的差异，支持忽略指定路径
func (s *jsonDiffServiceImpl) CompareJSONWithIgnore(json1, json2 string, ignorePaths []string) (JSONDiffResult, error) {
	var obj1, obj2 interface{}
	var result JSONDiffResult
	result.Changed = make(map[string]string)

	// 解析第一个JSON字符串
	if err := json.Unmarshal([]byte(json1), &obj1); err != nil {
		return result, fmt.Errorf("解析第一个JSON失败: %v", err)
	}

	// 解析第二个JSON字符串
	if err := json.Unmarshal([]byte(json2), &obj2); err != nil {
		return result, fmt.Errorf("解析第二个JSON失败: %v", err)
	}

	// 比较两个对象
	compareValues("", obj1, obj2, &result, ignorePaths)

	return result, nil
}

// compareValues 递归比较两个值并记录差异
func compareValues(path string, v1, v2 interface{}, result *JSONDiffResult, ignorePaths []string) {
	// 检查当前路径是否应该被忽略
	if shouldIgnorePath(path, ignorePaths) {
		return
	}

	// 处理nil值
	if v1 == nil && v2 == nil {
		return
	}
	// 检查null到非null的变更
	if v1 == nil {
		result.Changed[path] = fmt.Sprintf("值变更: null -> %v", v2)
		return
	}
	// 检查非null到null的变更
	if v2 == nil {
		result.Changed[path] = fmt.Sprintf("值变更: %v -> null", v1)
		return
	}

	// 获取值的类型
	type1 := reflect.TypeOf(v1)
	type2 := reflect.TypeOf(v2)

	// 如果类型不同
	if type1.Kind() != type2.Kind() {
		result.Changed[path] = fmt.Sprintf("类型变更: %v -> %v", type1.Kind(), type2.Kind())
		return
	}

	// 根据类型进行不同的比较
	switch t := v1.(type) {
	case map[string]interface{}:
		// 比较对象
		m2 := v2.(map[string]interface{})

		// 检查第一个对象中存在但第二个对象中不存在的键
		for k, v := range t {
			fullPath := buildPath(path, k)
			if _, exists := m2[k]; !exists {
				if !shouldIgnorePath(fullPath, ignorePaths) {
					result.Removed = append(result.Removed, fullPath)
				}
			} else {
				// 递归比较值
				compareValues(fullPath, v, m2[k], result, ignorePaths)
			}
		}

		// 检查第二个对象中存在但第一个对象中不存在的键
		for k := range m2 {
			if _, exists := t[k]; !exists {
				fullPath := buildPath(path, k)
				if !shouldIgnorePath(fullPath, ignorePaths) {
					result.Added = append(result.Added, fullPath)
				}
			}
		}

	case []interface{}:
		// 比较数组
		a2 := v2.([]interface{})

		// 如果数组长度不同，记录长度变更
		if len(t) != len(a2) {
			result.Changed[path] = fmt.Sprintf("数组长度变更: %d -> %d", len(t), len(a2))
		}

		// 比较对应位置的元素（比较到较短数组的长度）
		minLen := len(t)
		if len(a2) < minLen {
			minLen = len(a2)
		}

		for i := 0; i < minLen; i++ {
			compareValues(fmt.Sprintf("%s[%d]", path, i), t[i], a2[i], result, ignorePaths)
		}

		// 处理数组长度不同的情况
		if len(t) > len(a2) {
			// 第一个数组更长，标记多余元素为移除
			for i := len(a2); i < len(t); i++ {
				arrayPath := fmt.Sprintf("%s[%d]", path, i)
				if !shouldIgnorePath(arrayPath, ignorePaths) {
					result.Removed = append(result.Removed, arrayPath)
				}
			}
		} else if len(t) < len(a2) {
			// 第二个数组更长，标记多余元素为新增
			for i := len(t); i < len(a2); i++ {
				arrayPath := fmt.Sprintf("%s[%d]", path, i)
				if !shouldIgnorePath(arrayPath, ignorePaths) {
					result.Added = append(result.Added, arrayPath)
				}
			}
		}

	default:
		// 比较基本类型值
		if !reflect.DeepEqual(v1, v2) {
			result.Changed[path] = fmt.Sprintf("值变更: %v -> %v", v1, v2)
		}
	}
}

// buildPath 构建完整的键路径
func buildPath(parentPath, key string) string {
	if parentPath == "" {
		return key
	}
	return fmt.Sprintf("%s.%s", parentPath, key)
}

// shouldIgnorePath 检查给定路径是否应该被忽略，支持数组通配符 [*]
func shouldIgnorePath(path string, ignorePaths []string) bool {
	if ignorePaths == nil {
		return false
	}

	for _, ignorePath := range ignorePaths {
		if matchesPath(path, ignorePath) {
			return true
		}
	}
	return false
}

// matchesPath 检查路径是否匹配忽略模式，支持数组通配符
func matchesPath(path, pattern string) bool {
	// 精确匹配
	if path == pattern {
		return true
	}

	// 如果模式不包含通配符，则只进行精确匹配
	if !strings.Contains(pattern, "[*]") {
		return false
	}

	// 转换通配符模式为正则表达式
	// 将 [*] 替换为 \[\d+\] 来匹配任意数字索引
	regexPattern := regexp.QuoteMeta(pattern)
	regexPattern = strings.ReplaceAll(regexPattern, `\[\*\]`, `\[\d+\]`)
	regexPattern = "^" + regexPattern + "$"

	matched, err := regexp.MatchString(regexPattern, path)
	if err != nil {
		return false
	}

	return matched
}
