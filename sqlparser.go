package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)


func main() {
	//var converted_select select_clause

	argsWithProg := os.Args
	//Example SQL
 	//str1 := `select  sum(campo1), campo2, (select field, ffiend from (select test1, test2, test3 from internaltab1)) from (SELECT * FROM TAB3)
 //where t1 = 1 
 	//`
	 str1 := `select  sum(campo1), campo2 from table1, table2 where t1 = 'TEST STRING' `
//str1 := `select  campo1, campo2, (select field, ffiend from tab2) from tabela1 where t1 = 1`

	if argsWithProg != nil {
		if len(argsWithProg) > 1 && argsWithProg[1] != "" {
			str1 = argsWithProg[1]
		}
		//select campo1, campo2, (select field, ffiend from tab2) from tabela1 where t1 = 1
	}

	/*
		check_query_type(str1, &converted_select)

		fmt.Println("----------------------------------------------")
		fmt.Println("Resulting Object")
		fmt.Println("----------------------------------------------")
		fmt.Println(converted_select)
		fmt.Println("----------------------------------------------")
		fmt.Println("Original Query")
		fmt.Println("----------------------------------------------")
		fmt.Println(str1)
	*/
	var action ActionExec = InternalActionExec{}
	SetAction(action)
	result := tokenize_command(str1)
	execute_parsing_process(str1)
	fmt.Println(result)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type expression_unit struct {
	Index int
	Expression  string
}

type CommandTree struct {
	ClauseName   string
	TypeToken    string
	Clause		  string
	FullCommand  string
	CommandParts []CommandTree

}

type ActionExec interface {
	ExecAction(tree CommandTree)
}

type InternalActionExec struct{
	ActionExec
}



var _command_syntax_tree CommandTree
var _expressions []expression_unit
var _action ActionExec

func SetAction(action ActionExec) {
	_action = action
}

func (internalExec InternalActionExec) ExecAction(tree CommandTree){
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("CorrespondingAction")
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(tree)
}

func IndexStringSlice(slice [] string, value string) int {
    for p, v := range slice {
        if (v == value) {
            return p
        }
    }
    return -1
}

func control_hierarchy(expression string, opening_char string, ending_char string) string {
	result_expression := ""
	counter_hierarchy := 0 //send the expression with openingchar, please
	overall_counter := 0

	for _, element := range expression  {
		if string(element) == opening_char {
			counter_hierarchy++
		} else {
			if string(element) == ending_char {
				counter_hierarchy--
				if counter_hierarchy == 0 {
					result_expression = expression[strings.Index(expression, opening_char) + 1:overall_counter]
					break
				}
			}
		}
		overall_counter++
	}

	return result_expression
}

func control_string(expression string, opening_char string, ending_char string) string {
	result_expression := ""
	counter_hierarchy := 0 //send the expression with openingchar, please
	overall_counter := 0

	for _, element := range expression  {
		if string(element) == opening_char && counter_hierarchy == 0 {
			counter_hierarchy++
		} else {
			if string(element) == ending_char {
				result_expression = expression[strings.Index(expression, opening_char) + 1:overall_counter]
				break
			}
		}
		
		overall_counter++
	}

	return result_expression
}


func control_hierarchy_tokenized(expression [] string, opening_char string, ending_char string) []string {
	var result_expression []string
	counter_hierarchy := 0 //send the expression with openingchar, please
	overall_counter := 0

	for _, element := range expression {
		if string(element) == opening_char {
			counter_hierarchy++
		} else if string(element) == ending_char {
			counter_hierarchy--
			if counter_hierarchy == 0 {
				result_expression = expression[IndexStringSlice(expression, opening_char) + 1:overall_counter]
				break
			}
		}
		overall_counter++
	}
	if len(result_expression) < 1{
		result_expression = expression[IndexStringSlice(expression, opening_char) + 1:]
	}

	return result_expression
}

var count int = 0
func get_all_sub_expressions(current_index int){
	for strings.Index(_expressions[current_index].Expression, "'") > -1 && count < 15 {
		sub_expresion := control_string(_expressions[current_index].Expression, "'", "'")
		fmt.Println("SUB STRING-------------------------------")
		fmt.Println(sub_expresion)
		fmt.Println("-----------------------------------------")

		unit := expression_unit{Index:len(_expressions), Expression:"'" + sub_expresion + "'"}
		strindex := " {" + strconv.Itoa(len(_expressions)) + "} "
		_expressions[current_index].Expression = strings.Replace(_expressions[current_index].Expression, "'" + sub_expresion + "'", strindex, 1)
		_expressions = append(_expressions, unit)
		count++
	}

	for strings.Index(_expressions[current_index].Expression, "(") > -1{
		sub_expresion := control_hierarchy(_expressions[current_index].Expression, "(", ")")
		unit := expression_unit{Index:len(_expressions), Expression:sub_expresion}
		strindex := " {" + strconv.Itoa(len(_expressions)) + "} "
		_expressions[current_index].Expression = strings.Replace(_expressions[current_index].Expression, "(" + sub_expresion + ")", strindex, 1)
		_expressions = append(_expressions, unit)
		get_all_sub_expressions(unit.Index)
	}
}

func execute_parsing_process(command string) {
	command = strings.ToLower(command)
	unit := expression_unit{Index:0, Expression:command}
	_expressions = append(_expressions, unit)
	get_all_sub_expressions(0)
	fmt.Println("Expression-------------------------------")
	fmt.Println(_expressions)
	fmt.Println("-----------------------------------------")
	start_syntax_tree(command)
	// token_list := tokenize_command(command)
	// for token_index, _ := range token_list {
	// 	parse_token(&token_list, token_index)
	// }
}
func tokenize_command(command string) []string {
	re := regexp.MustCompile(`\S+`)
	submatchall := re.FindAllString(command, -1)
	return submatchall
}

func start_syntax_tree(command string){
	_command_syntax_tree = CommandTree{ClauseName:"master",
		TypeToken:"master",
		FullCommand:command}
	first_exmpression := strings.Trim(_expressions[0].Expression, " ")
	switch strings.ToLower(strings.Trim(first_exmpression[0:strings.Index(first_exmpression, " ")], " ")){
	case "select":
		fmt.Println("start_syntax_tree-------------------------------")
		fmt.Println("SELECT")
		fmt.Println("------------------------------------------------")
		_command_syntax_tree.CommandParts = append(_command_syntax_tree.CommandParts, CommandTree{ClauseName:"SELECT",
		TypeToken:"SELECT"})
		//, FullCommand:_expressions[0].Expression })

		parse_select_regions(_expressions[0].Expression, &_command_syntax_tree.CommandParts[len(_command_syntax_tree.CommandParts)-1])

		fmt.Println("End Syntax Tree-------------------------------")
		fmt.Println(_command_syntax_tree)
		fmt.Println("------------------------------------------------")
		//METHOD to break down select

		break;
	default:
		fmt.Println("start_syntax_tree-------------------------------")
		fmt.Println(first_exmpression)
		fmt.Println("------------------------------------------------")
		break;
	}
}

func parse_select_regions(expression string, tree * CommandTree){
	tokens := tokenize_command(expression)

	fmt.Println("TOKENS-------------------------------")
	fmt.Println(tokens)
	fmt.Println("------------------------------------------------")


	if (IndexStringSlice(tokens, "select")> -1){
		tokenized_fields := control_hierarchy_tokenized(tokens, "select", "from")
		fmt.Println("tokenized_fields-------------------------------")
		fmt.Println(tokenized_fields)
		fmt.Println("------------------------------------------------")
		tree_part := CommandTree{ClauseName:"select",
		TypeToken:"FIELDS_SELECT",
		FullCommand:expression}
		tree.CommandParts = append(tree.CommandParts, tree_part) // has to call get_tokens_as_tree
		get_tokens_as_tree(tokenized_fields, &tree.CommandParts[len(tree.CommandParts)-1])
	}
	if (IndexStringSlice(tokens, "from")> -1){
		tokenized_tables := control_hierarchy_tokenized(tokens, "from", "where")
		fmt.Println("tokenized_tables-------------------------------")
		fmt.Println(tokenized_tables)
		fmt.Println("------------------------------------------------")
		tree_part := CommandTree{ClauseName:"from",
		TypeToken:"tables_from"}//,FullCommand:expression}
		tree.CommandParts = append(tree.CommandParts, tree_part) // has to call get_tokens_as_tree
		get_tokens_as_tree(tokenized_tables, &tree.CommandParts[len(tree.CommandParts)-1])
	}
	if (IndexStringSlice(tokens, "where")> -1){
		tokenized_filters := control_hierarchy_tokenized(tokens, "where", "go;")
		fmt.Println("tokenized_filters-------------------------------")
		fmt.Println(tokenized_filters)
		fmt.Println("------------------------------------------------")
		tree_part := CommandTree{ClauseName:"where",
		TypeToken:"where_fields"}//, FullCommand:expression}
		tree.CommandParts = append(tree.CommandParts, tree_part) // has to call get_tokens_as_tree
		get_tokens_as_tree(tokenized_filters, &tree.CommandParts[len(tree.CommandParts)-1])
	}
}

func get_tokens_as_tree(tokenized_command []string, tree * CommandTree) []CommandTree{
	//var tree_curren []CommandTree
	for _, token := range tokenized_command {
		tree_part := get_command(token, tree)
		tree.CommandParts = append(tree.CommandParts, tree_part)
	}
	return tree.CommandParts
}


func get_command(command string, tree * CommandTree) CommandTree{
	//var tree CommandTree
	token := strings.Replace(command, ",", "", 1)//replace all maybe
	index_token := check_index(token)
	if  index_token > -1 {
		fmt.Println("get_command INDEX-------------------------------")
		fmt.Println(index_token)
		fmt.Println("------------------------------------------------")
		//I'll need a recursion here to parse_select_regions, meaning also a pointer to the tree 
		if (strings.Index(_expressions[index_token].Expression, "'") > -1){
			//Treat as string
			tree = &CommandTree{ClauseName:_expressions[index_token].Expression, TypeToken:"STRING", Clause:_expressions[index_token].Expression}
		}else{
			parse_select_regions(_expressions[index_token].Expression, tree)
		}
	}else if _, err := strconv.Atoi(token); err == nil {
		tree = &CommandTree{ClauseName:token, TypeToken:"NUMBER", Clause:token}
	}else{
		switch token {
			case "select":
				tree = &CommandTree{ClauseName:"SELECT", TypeToken:"FIELDS_SELECT", Clause:token}
				//, FullCommand:command}
				break
			case "inner":
				tree = &CommandTree{ClauseName:"INNER", TypeToken:"JOIN_TYPE", Clause:token}
				//, FullCommand:command}
				break
			case "left":
				tree = &CommandTree{ClauseName:"LEFT", TypeToken:"JOIN_TYPE", Clause:token}
				//, FullCommand:command}
				break
			case "right":
				tree = &CommandTree{ClauseName:"RIGHT", TypeToken:"JOIN_TYPE", Clause:token}
				//, FullCommand:command}
				break
			case "outer":
				tree = &CommandTree{ClauseName:"OUTER", TypeToken:"JOIN_TYPE", Clause:token}
				//, FullCommand:command}
				break
			case "join":
				tree = &CommandTree{ClauseName:"JOIN", TypeToken:"JOIN", Clause:token}
				//, FullCommand:command}
				break
			case "sum":
				tree = &CommandTree{ClauseName:"SUM", TypeToken:"RESERVED_FUNCTION", Clause:token}
				//, FullCommand:command}
				break
			case "group":
				tree = &CommandTree{ClauseName:"GROUP", TypeToken:"RESERVED_FUNCTION", Clause:token}
				//, FullCommand:command}
				break
			case "max":
				tree = &CommandTree{ClauseName:"MAX", TypeToken:"RESERVED_FUNCTION", Clause:token}
				//, FullCommand:command}
				break
			case "distinct":
				tree = &CommandTree{ClauseName:"DISTINCT", TypeToken:"RESERVED_FUNCTION", Clause:token}
				//, FullCommand:command}
				break
			case "=":
				tree = &CommandTree{ClauseName:"EQUALS", TypeToken:"OPERATOR", Clause:token}
				//, FullCommand:command}
				break
			case ">":
				tree = &CommandTree{ClauseName:"BIGGER_THAN", TypeToken:"OPERATOR", Clause:token}
				//, FullCommand:command}
				break
			case "<":
				tree = &CommandTree{ClauseName:"SMALLER_THAN", TypeToken:"OPERATOR", Clause:token}
				//, FullCommand:command}
				break
			case "and":
				tree = &CommandTree{ClauseName:"AND", TypeToken:"OPERATOR", Clause:token}
				//, FullCommand:command}
				break
			case "or":
				tree = &CommandTree{ClauseName:"OR", TypeToken:"OPERATOR", Clause:token}
				//, FullCommand:command}
				break
			default:
				if (tree.ClauseName == "select"){
					tree = &CommandTree{ClauseName:"SELECT", TypeToken:"FIELD_SELECT_TO_SHOW", Clause:token}
					//, FullCommand:command}
				}
				if (tree.ClauseName == "from"){
					tree = &CommandTree{ClauseName:"FROM", TypeToken:"TABLE_FROM", Clause:token}
					//, FullCommand:command}
				}
				if (tree.ClauseName == "where"){
					tree = &CommandTree{ClauseName:"WHERE", TypeToken:"FIELD_FILTER", Clause:token}
					//, FullCommand:command}
				}
				break
		}
	}

	_action.ExecAction(*tree)
	return *tree
}



func check_index(command string) int {
	re := regexp.MustCompile(`{\d+}`)
	submatchall := re.FindAllString(command, -1)
	result := -1
	if len(submatchall) > 0{
		number_string := strings.Replace(strings.Replace(submatchall[0], "{", "", 1), "}", "", 1)
		result, _ = strconv.Atoi(number_string)
	}
	return result
}