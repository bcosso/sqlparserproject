package sqlparserproject

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)


///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type expression_unit struct {
	Index      int
	Expression string
}

type ExpressionPart struct {
	BeginIdentifier string
	EndIdentifier   string
}
type ExpressionTemplate struct {
	ExpressionType string
	Parts          []ExpressionPart
}

type CommandTree struct {
	ClauseName   string
	TypeToken    string
	Clause       string
	Alias		 string
	Prefix		 string
	FullCommand  string
	CommandParts []CommandTree
}


type LogicGates struct {
	CommandLeft  interface{}
	Gate    	 string
	CommandRight interface{}
	ChildGates 	 []LogicGates
}


type ActionExec interface {
	ExecAction(tree * CommandTree)
	ExecActionFinal(tree CommandTree)
}

type InternalActionExec struct {
	ActionExec
}

var _command_syntax_tree CommandTree
var _expressions []expression_unit
var _action ActionExec

func isInteger(val float64) bool {
    return val == float64(int(val))
}

func SetAction(action ActionExec) {
	_action = action
}

func (internalExec InternalActionExec) ExecAction(tree * CommandTree) {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("CorrespondingAction")
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(tree)
}

func (internalExec InternalActionExec) ExecActionFinal(tree CommandTree) {
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("CorrespondingFinalAction")
	fmt.Println("-----------------------------------------------------------")
	fmt.Println(tree)
}

func IndexStringSlice(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

func control_hierarchy(expression string, opening_char string, ending_char string) string {
	result_expression := ""
	counter_hierarchy := 0 //send the expression with openingchar, please
	overall_counter := 0

	for _, element := range expression {
		if string(element) == opening_char {
			counter_hierarchy++
		} else {
			if string(element) == ending_char {
				counter_hierarchy--
				if counter_hierarchy == 0 {
					result_expression = expression[strings.Index(expression, opening_char)+1 : overall_counter]
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

	for _, element := range expression {
		if string(element) == opening_char && counter_hierarchy == 0 {
			counter_hierarchy++
		} else {
			if string(element) == ending_char {
				result_expression = expression[strings.Index(expression, opening_char)+1 : overall_counter]
				break
			}
		}

		overall_counter++
	}

	return result_expression
}

func control_hierarchy_tokenized(expression []string, opening_char string, ending_char string) []string {
	var result_expression []string
	counter_hierarchy := 0 //send the expression with openingchar, please
	overall_counter := 0

	for _, element := range expression {
		if string(element) == opening_char {
			counter_hierarchy++
		} else if string(element) == ending_char {
			counter_hierarchy--
			if counter_hierarchy == 0 {
				result_expression = expression[IndexStringSlice(expression, opening_char)+1 : overall_counter]
				break
			}
		}
		overall_counter++
	}
	if len(result_expression) < 1 {
		result_expression = expression[IndexStringSlice(expression, opening_char)+1:]
	}

	return result_expression
}

var count int = 0

func get_all_sub_expressions(current_index int) {
	for strings.Index(_expressions[current_index].Expression, "'") > -1 && count < 15 {
		sub_expresion := control_string(_expressions[current_index].Expression, "'", "'")
		unit := expression_unit{Index: len(_expressions), Expression: "'" + sub_expresion + "'"}
		strindex := " {" + strconv.Itoa(len(_expressions)) + "} "
		_expressions[current_index].Expression = strings.Replace(_expressions[current_index].Expression, "'"+sub_expresion+"'", strindex, 1)
		_expressions = append(_expressions, unit)
		count++
	}

	for strings.Index(_expressions[current_index].Expression, "(") > -1 {
		sub_expresion := control_hierarchy(_expressions[current_index].Expression, "(", ")")
		unit := expression_unit{Index: len(_expressions), Expression: sub_expresion}
		strindex := " {" + strconv.Itoa(len(_expressions)) + "} "
		_expressions[current_index].Expression = strings.Replace(_expressions[current_index].Expression, "("+sub_expresion+")", strindex, 1)
		_expressions = append(_expressions, unit)
		get_all_sub_expressions(unit.Index)
	}
}

func Execute_parsing_process(command string) CommandTree {
	_command_syntax_tree.CommandParts = nil
	_expressions = nil
	command = strings.ToLower(command)
	unit := expression_unit{Index: 0, Expression: command}
	_expressions = append(_expressions, unit)
	get_all_sub_expressions(0)
	start_syntax_tree(command)
	_action.ExecActionFinal(_command_syntax_tree)
	
	return _command_syntax_tree
}

func tokenize_command(command string) []string {
	re := regexp.MustCompile(`\S+`)
	submatchall := re.FindAllString(command, -1)
	return submatchall
}

func start_syntax_tree(command string) {
	_command_syntax_tree = CommandTree{ClauseName: "master",
		TypeToken:   "master",
		FullCommand: command}
	first_exmpression := strings.Trim(_expressions[0].Expression, " ")
	switch strings.ToLower(strings.Trim(first_exmpression[0:strings.Index(first_exmpression, " ")], " ")) {
	case "select":
		_command_syntax_tree.CommandParts = append(_command_syntax_tree.CommandParts, CommandTree{ClauseName: "SELECT_COMMAND",
			TypeToken: "SELECT"})
		parse_select_regions(_expressions[0].Expression, &_command_syntax_tree.CommandParts[len(_command_syntax_tree.CommandParts)-1])

		fmt.Println("End Syntax Tree-------------------------------")
		fmt.Println(_command_syntax_tree)
		fmt.Println("------------------------------------------------")
		break

	case "insert":
		_command_syntax_tree.CommandParts = append(_command_syntax_tree.CommandParts, CommandTree{ClauseName: "INSERT_COMMAND",
			TypeToken: "INSERT"})
		parse_insert_regions(_expressions[0].Expression, &_command_syntax_tree.CommandParts[len(_command_syntax_tree.CommandParts)-1])

		fmt.Println("End Syntax Tree-------------------------------")
		fmt.Println(_command_syntax_tree)
		fmt.Println("------------------------------------------------")
		break
	default:
		// fmt.Println("start_syntax_tree-------------------------------")
		// fmt.Println(first_exmpression)
		// fmt.Println("------------------------------------------------")
		break
	}
}

func parse_select_regions(expression string, tree *CommandTree) {
	tokens := tokenize_command(expression)

	if IndexStringSlice(tokens, "select") > -1 {
		tokenized_fields := control_hierarchy_tokenized(tokens, "select", "from")
		tree_part := CommandTree{ClauseName: "select", TypeToken: "FIELDS_SELECT", FullCommand: expression}
		tree.CommandParts = append(tree.CommandParts, tree_part) // has to call get_tokens_as_tree
		get_tokens_as_tree(tokenized_fields, &tree.CommandParts[len(tree.CommandParts)-1])
	}
	if IndexStringSlice(tokens, "from") > -1 {
		tokenized_tables := control_hierarchy_tokenized(tokens, "from", "where")
		tree_part := CommandTree{ClauseName: "from", TypeToken: "tables_from"} //,FullCommand:expression}
		tree.CommandParts = append(tree.CommandParts, tree_part) // has to call get_tokens_as_tree
		get_tokens_as_tree(tokenized_tables, &tree.CommandParts[len(tree.CommandParts)-1])
	}
	if IndexStringSlice(tokens, "where") > -1 {
		tokenized_filters := control_hierarchy_tokenized(tokens, "where", "go;")
		tree_part := CommandTree{ClauseName: "where", TypeToken: "where_fields"}
		tree.CommandParts = append(tree.CommandParts, tree_part) // has to call get_tokens_as_tree
		get_tokens_as_tree(tokenized_filters, &tree.CommandParts[len(tree.CommandParts)-1])
	}
}

func parse_insert_regions(expression string, tree *CommandTree) {
	tokens := tokenize_command(expression)

	if IndexStringSlice(tokens, "insert") > -1 {
		tokenized_fields := control_hierarchy_tokenized(tokens, "insert", "values")
		tree_part := CommandTree{ClauseName: "insert", TypeToken: "ADDRESSING_INSERT", FullCommand: expression}
		tree.CommandParts = append(tree.CommandParts, tree_part) // has to call get_tokens_as_tree
		get_tokens_as_tree(tokenized_fields, &tree.CommandParts[len(tree.CommandParts)-1])
	}
	if IndexStringSlice(tokens, "values") > -1 {
		tokenized_tables := control_hierarchy_tokenized(tokens, "values", ")")
		tree_part := CommandTree{ClauseName: "values", TypeToken: "VALUES_INSERT"} 
		tree.CommandParts = append(tree.CommandParts, tree_part) // has to call get_tokens_as_tree
		get_tokens_as_tree(tokenized_tables, &tree.CommandParts[len(tree.CommandParts)-1])
	}
}

func get_tokens_as_tree(tokenized_command []string, tree *CommandTree) []CommandTree {
	//var tree_curren []CommandTree
	index_token := 0 

	for index_token < len(tokenized_command) {
		tree_part := get_command(tokenized_command[index_token], tree, tokenized_command, &index_token)
		tree.CommandParts = append(tree.CommandParts, tree_part)
		index_token ++
	}
	return tree.CommandParts
}

func check_expression_containing_token(index_token int) string {
    index_expressions := index_token
	result_type := ""

	for index_expressions > 0 {
		
		if (strings.Index(strings.ToLower(_expressions[index_expressions].Expression), "select") == 0) || 
		(strings.Index(strings.ToLower(_expressions[index_expressions].Expression), "from") == 0) || 
		(strings.Index(strings.ToLower(_expressions[index_expressions].Expression), "where") == 0 ) || 
		(strings.Index(strings.ToLower(_expressions[index_expressions].Expression), "insert") == 0 ) || 
		(strings.Index(strings.ToLower(_expressions[index_expressions].Expression), "into") == 0 ) || 
		(strings.Index(strings.ToLower(_expressions[index_expressions].Expression), "update") == 0 ) || 
		(strings.Index(strings.ToLower(_expressions[index_expressions].Expression), "set") == 0) {

			result_type = strings.Trim(_expressions[index_expressions].Expression[0:strings.Index(strings.TrimLeft(_expressions[index_expressions].Expression, " "), " ")], " ")
			break
		} 
		index_expressions --
	}
	return result_type
}

func get_command(command string, tree *CommandTree, tokenized_command []string, index_tokenized_command * int) CommandTree {
	//var tree CommandTree
	token := strings.Replace(command, ",", "", 1) //replace all maybe
	index_token := check_index(token)
	
	if index_token > -1 {

		if strings.Index(_expressions[index_token].Expression, "'") > -1 {
			tree = &CommandTree{ClauseName: _expressions[index_token].Expression, TypeToken: "STRING", Clause: _expressions[index_token].Expression}
		} else if (*index_tokenized_command > 1 && strings.Index(tokenized_command[*index_tokenized_command - 2], "into") > -1) {
			tree = &CommandTree{ClauseName: _expressions[index_token].Expression, TypeToken: "COLUMNS", Clause: _expressions[index_token].Expression}
			get_tokens_as_tree(tokenize_command(_expressions[index_token].Expression),tree)
		} else if (strings.Index(strings.ToLower(_expressions[index_token].Expression), "select") == 0){
			tree = &CommandTree{ClauseName: "FIELDS", TypeToken: "FIELDS", Clause: _expressions[index_token].Expression}
			parse_select_regions(_expressions[index_token].Expression, tree)

		}else{
			tree = &CommandTree{ClauseName: "FIELDS", TypeToken: "FIELDS", Clause: _expressions[index_token].Expression}
			get_tokens_as_tree(tokenize_command(_expressions[index_token].Expression),tree)
		}
	} else if fNumber, err := strconv.ParseFloat(token, 64); err == nil {
		if isInteger(fNumber) {
			tree = &CommandTree{ClauseName: token, TypeToken: "INT", Clause: token}
		} else {
		tree = &CommandTree{ClauseName: token, TypeToken: "FLOAT64", Clause: token}
		}
	} else {
		switch token {
		case "select":
			tree = &CommandTree{ClauseName: "SELECT", TypeToken: "SELECT_COMMAND", Clause: token}
			break
		case "into":
			tree = &CommandTree{ClauseName: "INTO", TypeToken: "INTO_COMMAND", Clause: token}
			break
		case "in":
			tree = &CommandTree{ClauseName: "IN", TypeToken: "OPERATOR", Clause: token}
			break
		case "on":
			tree = &CommandTree{ClauseName: "ON", TypeToken: "ON_COMMAND", Clause: token}
			break
		case "inner":
			tree = &CommandTree{ClauseName: "INNER", TypeToken: "JOIN_TYPE", Clause: token}
			break
		case "left":
			tree = &CommandTree{ClauseName: "LEFT", TypeToken: "JOIN_TYPE", Clause: token}

			break
		case "right":
			tree = &CommandTree{ClauseName: "RIGHT", TypeToken: "JOIN_TYPE", Clause: token}

			break
		case "outer":
			tree = &CommandTree{ClauseName: "OUTER", TypeToken: "JOIN_TYPE", Clause: token}

			break
		case "join":
			tree = &CommandTree{ClauseName: "JOIN", TypeToken: "JOIN", Clause: token}

			break
		case "sum":
			tree = &CommandTree{ClauseName: "SUM", TypeToken: "RESERVED_FUNCTION", Clause: token}

			break
		case "group":
			tree = &CommandTree{ClauseName: "GROUP", TypeToken: "RESERVED_FUNCTION", Clause: token}

			break
		case "max":
			tree = &CommandTree{ClauseName: "MAX", TypeToken: "RESERVED_FUNCTION", Clause: token}

			break
		case "distinct":
			tree = &CommandTree{ClauseName: "DISTINCT", TypeToken: "RESERVED_FUNCTION", Clause: token}

			break
		case "=":
			tree = &CommandTree{ClauseName: "EQUALS", TypeToken: "OPERATOR", Clause: token}

			break
		case "*=":
			tree = &CommandTree{ClauseName: "LEFT_JOIN", TypeToken: "OPERATOR", Clause: token}

			break
		case "=*":
			tree = &CommandTree{ClauseName: "RIGHT_JOIN", TypeToken: "OPERATOR", Clause: token}

			break
		case ">":
			tree = &CommandTree{ClauseName: "BIGGER_THAN", TypeToken: "OPERATOR", Clause: token}

			break
		case "<":
			tree = &CommandTree{ClauseName: "SMALLER_THAN", TypeToken: "OPERATOR", Clause: token}

			break
		case "and":
			tree = &CommandTree{ClauseName: "AND", TypeToken: "OPERATOR", Clause: token}

			break
		case "or":
			tree = &CommandTree{ClauseName: "OR", TypeToken: "OPERATOR", Clause: token}

			break
		case "as":
			if len(tokenized_command) > (*index_tokenized_command) + 1{
				
					treeChild := &((*tree).CommandParts[len(tree.CommandParts) -1])
					(*treeChild).Alias = tokenized_command[(*index_tokenized_command) + 1]
					*index_tokenized_command += 1
			}
			break

		default:
			if tree.ClauseName == "select" {
				tree = &CommandTree{ClauseName: "SELECT", TypeToken: "FIELD_SELECT_TO_SHOW", Clause: token}
				CheckForPrefixes(tree)
	
			} else
			if tree.ClauseName == "from" {
				tree = &CommandTree{ClauseName: "FROM", TypeToken: "TABLE_FROM_COMMAND", Clause: token}
	
			} else
			if tree.ClauseName == "where" {
				tree = &CommandTree{ClauseName: "WHERE", TypeToken: "FIELD_FILTER", Clause: token}
				CheckForPrefixes(tree)
	
			} else
			if tree.ClauseName == "into" {
				tree = &CommandTree{ClauseName: "INTO", TypeToken: "TABLE_INTO_COMMAND", Clause: token}
	
			} else
			if tree.ClauseName == "values" {
				tree = &CommandTree{ClauseName: "VALUES", TypeToken: "COLUMN_VALUES_COMMAND", Clause: token}
	
			}else if (*index_tokenized_command > 0 && strings.Index(tokenized_command[*index_tokenized_command - 1], "into") > -1){
				tree = &CommandTree{ClauseName: "TABLE", TypeToken: "TABLE", Clause: token}
			} else {
				tree = &CommandTree{ClauseName: "FIELD", TypeToken: "FIELD", Clause: token}
			}
			break
		}		
	}
	if len(tokenized_command) > (*index_tokenized_command) + 2{
		if tokenized_command[(*index_tokenized_command) + 1] == "as" {
			(*tree).Alias = strings.Replace(tokenized_command[(*index_tokenized_command) + 2], ",", "", -1)
			*index_tokenized_command += 2
		}
	}

	_action.ExecAction(tree)
	return *tree
}

func CheckForPrefixes(tree *CommandTree ){
	dot := strings.Index((*tree).Clause, ".")
	if dot > -1{
		(*tree).Prefix = (*tree).Clause[:dot]
		(*tree).Clause = (*tree).Clause[dot+1:]
	}
}

func check_index(command string) int {
	re := regexp.MustCompile(`{\d+}`)
	submatchall := re.FindAllString(command, -1)
	result := -1
	if len(submatchall) > 0 {
		number_string := strings.Replace(strings.Replace(submatchall[0], "{", "", 1), "}", "", 1)
		result, _ = strconv.Atoi(number_string)
	}
	return result
}
