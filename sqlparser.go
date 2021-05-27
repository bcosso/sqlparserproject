package main

import (
	"fmt"
	"regexp"
	"strings"
)


type QueryKeywords struct {
	Keyword string
	Position int
}

type select_clause struct {
	clause_name string
	alias string
	value string
	operator string
	select_fields []select_clause
	table_list []select_clause
	where_list []select_clause
}

func main() {
	var converted_select select_clause
	str1 := "select campo1, campo2, (select field, ffiend from tab2) from tabela1 where t1 = 1"
	check_query_type(str1, &converted_select)

	fmt.Println("----------------------------------------------")
	fmt.Println("Resulting Object")
	fmt.Println("----------------------------------------------")	
	fmt.Println(converted_select)
	fmt.Println("----------------------------------------------")
	fmt.Println("Original Query")
	fmt.Println("----------------------------------------------")
	fmt.Println(str1)
}

func get_select_object(input_string string, result * select_clause){
	
	var result_select_fields []select_clause
	var result_select_from []select_clause
	var result_select_where []select_clause


	result_select_fields = get_select_fields(input_string, "select", " from")
	result_select_from = get_select_fields(input_string, "from", " where")
	result_select_where = get_select_fields(input_string, "where", "")

	result.select_fields = result_select_fields
	result.table_list = result_select_from
	result.where_list = result_select_where
}


func get_select_fields(input_string string, command string, endcommand string) []select_clause{
	var result []select_clause

	if strings.Index(input_string, command) < 0 { return nil }
	re := regexp.MustCompile(command + ` [.*]` + endcommand)
	fmt.Printf("Pattern: %v\n", re.String()) // print pattern

	submatchall := re.Split(input_string, -1)
	for _, element := range submatchall {
		if endcommand != "" {
			if strings.Index(element, endcommand) > 0 {
				element = element[0:strings.Index(element, endcommand)]
			}
		}
		if strings.Index(element, command) > 0 {
			element = element[strings.Index(element, command): len(element) -1]
		}
		var result_ind select_clause
		result_ind.clause_name = strings.Trim(element, command + " ")
		result = append(result, result_ind)
		fmt.Println("----------------------------------------------")
		fmt.Println("Element of query")
		fmt.Println("----------------------------------------------")
		fmt.Println(element)
	}

	return result
}


func get_sub_expression(expression string, opening_char string, ending_char string) string{
	result_expression := ""
	counter_hierarchy := 0 //send the expression with openingchar, please
	counter := 0

	for _, element := range expression {
		if string(element) == opening_char{
			counter_hierarchy ++
		}else if string(element) == ending_char {
			counter_hierarchy --
			if counter_hierarchy == 0 {
				result_expression = expression[0:counter]
				break
			} 
		}
		counter++
	}	

	return result_expression
}

func check_query_type(expression string, object_query * select_clause){
	if strings.Index(strings.ToLower(expression), "select") < 2 {//have to make a trim
		//after getting the fields from the first parenthesis, get latest keyword before the opening of parenthesus
		sub_expression := get_sub_expression(expression[strings.Index(expression, "("):len(expression)], "(", ")")
		if object_query == nil {
			object_query = new(select_clause)
		}
		
		if sub_expression != "" {
			latest_keyword_ocurrence_string := expression[0:strings.Index(expression, "(")]
			latest_keyword_list := find_location_expression(latest_keyword_ocurrence_string)
			fmt.Println(latest_keyword_list)		
			expression = strings.Replace(expression, sub_expression, "", 1)
			get_select_object(expression, object_query)
			fmt.Println(object_query)
			if len(latest_keyword_list) > 0{
				if latest_keyword_list[0].Keyword == "SELECT"{
					check_query_type(sub_expression, &object_query.select_fields[len(object_query.select_fields) - 1])
				}else if latest_keyword_list[0].Keyword == "FROM"{
					check_query_type(sub_expression, &object_query.table_list[len(object_query.table_list)  - 1])
				}else if latest_keyword_list[0].Keyword == "WHERE"{
					check_query_type(sub_expression, &object_query.where_list[len(object_query.where_list)  - 1])
				}
			}
		}else{
			get_select_object(expression, object_query)
			//return item
		}
	}
}


func find_location_expression(expression string)[]QueryKeywords {
	keyword_list := []QueryKeywords{
		{"SELECT", 0},
		{"FROM", 0},
		{"WHERE", 0},
	}
	
	var result []QueryKeywords
	countInt := 0
	for _, element := range keyword_list {
		currentIndex := strings.Index(strings.ToUpper(expression), element.Keyword)
		element.Position = currentIndex
		if countInt > 0 && element.Position > -1 && currentIndex < result[countInt-1].Position{
			result = append(result, result[countInt-1])
			result[countInt-1] = element
		}else{
			result = append(result, element)
		}
		countInt ++

	}
	fmt.Println("----------------------------------------------")
	fmt.Println("find_location_expression")
	fmt.Println("----------------------------------------------")
	fmt.Println(result)
	return result
}