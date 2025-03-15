package handlers

// import (
// 	"casino_loyalty_reward_system/internal/types"
// 	"net/http"
// 	"strconv"
// )

// func GetPaginationFilterFromRequestQueries(r *http.Request) types.PaginationFilter {
// 	var pagination PaginationFilter
// 	if query := r.FormValue("query"); query != "" {
// 		pagination.Query = &query
// 	}

// 	if offset := r.FormValue("offset"); offset != "" {
// 		c, err := strconv.Atoi(offset)
// 		if err != nil {
// 			return pagination
// 		}
// 		pagination.Offset = &c
// 	}

// 	if count := r.FormValue("count"); count != "" {
// 		c, err := strconv.Atoi(count)
// 		if err != nil {
// 			return pagination
// 		}
// 		pagination.Count = &c
// 	}

// 	if sortProperty := r.FormValue("sort_field"); sortProperty != "" {
// 		pagination.SortProperty = &sortProperty
// 	}

// 	if sortOrder := r.FormValue("sort_order"); sortOrder != "" {
// 		if sortOrder == "-1" {
// 			pagination.SortOrder = StringToPtr("DESC")
// 		}
// 		if sortOrder == "1" {
// 			pagination.SortOrder = StringToPtr("ASC")
// 		}
// 	}

// 	return pagination
// }

// func StringToPtr(value string) *string {
// 	return &value
// }

// func ToPostgresWhereAndArgs(prefix string, filter any) ([]any, []string) {
// 	var (
// 		arguments []any
// 		sqlWhere  []string
// 	)

// 	fieldsNum := reflect.TypeOf(filter).NumField()
// 	for i := 0; i < fieldsNum; i++ {
// 		field := reflect.ValueOf(filter).Field(i)

// 		if field.Kind() != reflect.Struct && field.IsValid() && !field.IsNil() {
// 			column := reflect.TypeOf(filter).Field(i).Tag.Get("schema")
// 			queryType := reflect.TypeOf(filter).Field(i).Tag.Get("compare")
// 			if column == "" || queryType == "" {
// 				continue
// 			}

// 			split := strings.Split(queryType, ",")
// 			arguments = append(arguments, reflect.ValueOf(filter).FieldByIndex([]int{i}).Interface())
// 			if split[0] == "global" {
// 				var globalWhere []string
// 				for c := 1; c < len(split); c++ {
// 					globalWhere = append(globalWhere, fmt.Sprintf("%s%s ILIKE '%%' || $%d || '%%'", prefix, split[c], len(arguments)))
// 				}
// 				sqlWhere = append(sqlWhere, fmt.Sprintf("(%s)", strings.Join(globalWhere, " OR ")))
// 			} else if queryType == "like" {
// 				sqlWhere = append(sqlWhere, fmt.Sprintf("%s%s ILIKE '%%' || $%d || '%%'", prefix, column, len(arguments)))
// 			} else if queryType == "eq" {
// 				sqlWhere = append(sqlWhere, fmt.Sprintf("%s%s = $%d", prefix, column, len(arguments)))
// 			}
// 		}
// 	}

// 	return arguments, sqlWhere
// }