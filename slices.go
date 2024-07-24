package butils

import "fmt"

// สำหรับเช็คค่า slice ว่ามี string ของเราอยู่ในนั้นหรือไม่
func Contains[V string | []string](s []string, val V) bool {
	if val, ok := any(val).(string); ok {
		for _, v := range s {
			if v == val {
				return true
			}
		}
	}

	if val, ok := any(val).([]string); ok {
		check := []string{}
		for _, str := range val {
			for _, v := range s {
				if v == str {
					check = append(check, str)
					continue
				}
			}
		}
		if len(check) == len(val) {
			return true
		}
	}

	return false
}

// สำหรับ Filter slice โดย filter จากค่า return ของ function ที่เป็นจริง ส่ง type ไหนมา จะส่ง type นั้นกลับไป
func Filter[T any](v []T, f func(T) bool) []T {
	result := make([]T, 0, len(v))
	for _, s := range v {
		if f(s) {
			result = append(result, s)
		}
	}
	return result
}

// สำหรับตรวจสอบค่าใน slice ว่ามีค่าที่ตรงตามเงื่อนไขของ function อย่างน้อย 1 ตัว หรือไม่
func Some[T any](v []T, f func(T) bool) bool {
	for _, s := range v {
		if f(s) {
			return true
		}
	}
	return false
}

// สำหรับตรวจสอบค่าใน slice ว่ามีค่าที่ตรงตามเงื่อนไขของ function ทุกตัว หรือไม่
func Every[T any](v []T, f func(T) bool) bool {
	for _, s := range v {
		if !f(s) {
			return false
		}
	}
	return true
}

// สำหรับสร้าง Slice ใหม่ โดยมีค่าเท่ากับผลลัพธ์ของ function ที่ส่งค่าแต่ละ element เข้าไป
func Map[T any, R any](v []T, f func(T) R) (result []R) {
	for _, s := range v {
		result = append(result, f(s))
	}
	return
}

// สำหรับตรวจสอบค่าใน slice ว่ามีค่าที่ตรงตามเงื่อนไขของ function อย่างน้อย 1 ตัว หรือไม่ จากนั้นจะส่งค่านั้นมาพร้อม index
func Find[T any](v []T, f func(T) bool) (found T, idx int, err error) {
	for idx, s := range v {
		if f(s) {
			return s, idx, nil
		}
	}
	return found, -1, fmt.Errorf("not found")
}

// สำหรับลบค่าที่ซ้ำกันใน slice ออก
func SetUnique[T comparable](v []T) (result []T) {
	seen := make(map[T]struct{})
	result = make([]T, 0)
	for _, s := range v {
		if _, ok := seen[s]; !ok {
			result = append(result, s)
			seen[s] = struct{}{}
		}
	}
	return
}
