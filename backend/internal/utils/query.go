package utils

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func ApplyPagination(query *gorm.DB, page, pageSize int) *gorm.DB {
	offset := (page - 1) * pageSize
	return query.Offset(offset).Limit(pageSize)
}

func ApplySorting(query *gorm.DB, sortBy, sortOrder string, allowedFields map[string]string) *gorm.DB {
	if sortBy == "" {
		return query.Order("created_at DESC")
	}

	column, ok := allowedFields[sortBy]
	if !ok {
		return query.Order("created_at DESC")
	}

	order := strings.ToUpper(sortOrder)
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	return query.Order(fmt.Sprintf("%s %s", column, order))
}

func ApplySearch(query *gorm.DB, search string, fields ...string) *gorm.DB {
	if search == "" || len(fields) == 0 {
		return query
	}

	pattern := "%" + strings.ToLower(search) + "%"
	conditions := make([]string, len(fields))
	args := make([]interface{}, len(fields))
	for i, field := range fields {
		conditions[i] = fmt.Sprintf("LOWER(%s) LIKE ?", field)
		args[i] = pattern
	}

	return query.Where(strings.Join(conditions, " OR "), args...)
}
