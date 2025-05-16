package generator

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/devOpifex/obfuscator/environment"
)

// Configuration for the R code generation
const (
	minVariables   = 1
	maxVariables   = 3
	minOperations  = 1
	maxOperations  = 5
	variablePrefix = "xXvArXx"
)

// R variable types
var variableTypes = []string{
	"numeric",
	"character",
	"logical",
	"vector",
	"data.frame",
}

// Base R functions only (no package dependencies)
var rFunctions = []string{
	"mean", "median", "sum", "min", "max",
	"length", "nrow", "ncol", "head", "tail",
	"sort", "rev", "unique", "sample", "paste",
	"substr", "toupper", "tolower", "round", "ceiling",
	"floor", "abs", "sqrt", "log", "exp",
}

// GenerateRandomRCode generates a random but valid base R script
// with proper variable tracking, no imports, and proper TRUE/FALSE
func Generate() string {
	var code strings.Builder

	// Map to track existing variables and their types
	variableMap := make(map[string]string)

	// Create variables
	numVars := rand.Intn(maxVariables-minVariables+1) + minVariables

	for i := 0; i < numVars; i++ {
		varName := fmt.Sprintf("%s%d", variablePrefix, i+1)
		varName = environment.Mask(varName)
		varType := variableTypes[rand.Intn(len(variableTypes))]

		// Add variable to map
		variableMap[varName] = varType

		code.WriteString(generateVariable(varName, varType))
		code.WriteString(";")
	}

	// Get slice of variable names for random selection
	variables := make([]string, 0, len(variableMap))
	for v := range variableMap {
		variables = append(variables, v)
	}

	// Generate operations
	numOps := rand.Intn(maxOperations-minOperations+1) + minOperations
	for i := 0; i < numOps; i++ {
		newVar, opCode := generateOperation(variables)
		if newVar != "" {
			// Track newly created variables
			variableMap[newVar] = environment.Mask("xXxXdErIvEdxX")
			variables = append(variables, newVar)
		}
		code.WriteString(opCode)
		code.WriteString(";")
	}

	// Add a final assignment without any printing
	finalVar := variables[rand.Intn(len(variables))]
	resultVar := environment.Mask("xXxXrEsUlT_xXxX")
	variableMap[resultVar] = environment.Mask("xXxXdErIvEdxX")
	fnName := environment.Mask("xXxXfN_xXxX")
	code.WriteString(fmt.Sprintf("%s=\\(){tryCatch(%s(%s),error=\\(e){NA})};%s=tryCatch(%s(),error=\\(e){NA});",
		fnName,
		rFunctions[rand.Intn(len(rFunctions))],
		finalVar,
		resultVar,
		fnName))

	return code.String()
}

// generateVariable creates a random variable of the specified type
func generateVariable(name string, varType string) string {
	switch varType {
	case "numeric":
		return fmt.Sprintf("%s=%.2f", name, rand.Float64()*100)
	case "character":
		words := []string{"apple", "banana", "cherry", "date", "elderberry", "cake"}
		for i := range words {
			words[i] = environment.Mask(words[i])
		}
		return fmt.Sprintf("%s=\"%s\"", name, words[rand.Intn(len(words))])
	case "logical":
		// Using proper R boolean capitalization
		if rand.Intn(2) == 1 {
			return fmt.Sprintf("%s=T", name)
		} else {
			return fmt.Sprintf("%s=F", name)
		}
	case "vector":
		length := rand.Intn(10) + 3
		vectorType := rand.Intn(3)
		switch vectorType {
		case 0: // numeric vector
			return fmt.Sprintf("%s=c(%s)", name, generateRandomVector(length))
		case 1: // character vector
			chars := []string{"a", "b", "c", "d", "e"}
			for i := range chars {
				chars[i] = environment.Mask(chars[i])
			}
			return fmt.Sprintf("%s=c(\"%s\")[1:%d]", name, strings.Join(chars, "\",\""), length)
		case 2: // logical vector
			return fmt.Sprintf("%s=c(%s)", name, generateRandomLogicalVector(length))
		}
	case "data.frame":
		return fmt.Sprintf("%s=data.frame("+environment.Mask("xXxxXsisxAasd")+"=c(%s),"+environment.Mask("xXxxAasd")+"=c(%s),"+environment.Mask("xXxxAaWdxXXsjj")+"=c(%s))",
			name,
			generateRandomVector(5),
			generateRandomVector(5),
			generateRandomLogicalVector(5))
	default:
		return fmt.Sprintf("%s=NA", name)
	}
	return fmt.Sprintf("%s=NULL", name)
}

// generateRandomVector creates a random numeric vector of specified length
func generateRandomVector(length int) string {
	elements := make([]string, length)
	for i := 0; i < length; i++ {
		elements[i] = fmt.Sprintf("%.1f", rand.Float64()*100)
	}
	return strings.Join(elements, ",")
}

// generateRandomLogicalVector creates a random logical vector with proper R syntax
func generateRandomLogicalVector(length int) string {
	elements := make([]string, length)
	for i := 0; i < length; i++ {
		if rand.Intn(2) == 1 {
			elements[i] = "T"
		} else {
			elements[i] = "F"
		}
	}
	return strings.Join(elements, ",")
}

// generateOperation creates a random operation using existing variables
// Returns the name of any new variable created and the operation code
func generateOperation(variables []string) (string, string) {
	if len(variables) == 0 {
		return "", "# No variables available for operations"
	}

	opType := rand.Intn(4)

	switch opType {
	case 0: // Assignment with function
		sourceVar := variables[rand.Intn(len(variables))]
		function := rFunctions[rand.Intn(len(rFunctions))]
		newVar := fmt.Sprintf("result_%s_%d", sourceVar, rand.Intn(1000))
		newVar = environment.Mask(newVar)
		return newVar, fmt.Sprintf("%s=tryCatch(%s(%s),error=\\(e){NA})", newVar, function, sourceVar)

	case 1: // Conditional operation
		targetVar := variables[rand.Intn(len(variables))]
		newVar := fmt.Sprintf("%s_modified", targetVar)
		newVar = environment.Mask(newVar)
		return newVar, fmt.Sprintf("if(is.numeric(%s)&&length(%s)>0){%s=%s^2;}else{%s=%s;}",
			targetVar, targetVar, newVar, targetVar, newVar, targetVar)

	case 2: // For loop with assignment
		targetVar := variables[rand.Intn(len(variables))]
		newVar := fmt.Sprintf("%s_processed", targetVar)
		newVar = environment.Mask(newVar)
		return newVar, fmt.Sprintf("%s=%s;for(i in 1:3){if(is.numeric(%s)){%s=%s+i;};}",
			newVar, targetVar, newVar, newVar, newVar)

	case 3: // Combining two variables with ifelse
		if len(variables) < 2 {
			// Fall back to modifying a single variable if we don't have at least 2
			targetVar := variables[rand.Intn(len(variables))]
			newVar := fmt.Sprintf("combined_%d", rand.Intn(1000))
			newVar = environment.Mask(newVar)
			return newVar, fmt.Sprintf("%s=ifelse(is.numeric(%s),%s*2,NA)",
				newVar, targetVar, targetVar)
		}

		// Select two different variables
		idx1 := rand.Intn(len(variables))
		idx2 := (idx1 + 1 + rand.Intn(len(variables)-1)) % len(variables)
		var1 := variables[idx1]
		var2 := variables[idx2]
		newVar := fmt.Sprintf("combined_%s_%s", var1, var2)
		newVar = environment.Mask(newVar)

		// Create more sophisticated operations combining variables
		ops := []string{
			fmt.Sprintf("%s=c(%s=%s,%s=%s)", newVar, var1, var1, var2, var2),
			fmt.Sprintf("%s=ifelse(is.numeric(%s)&&is.numeric(%s),%s+%s,NA)", newVar, var1, var2, var1, var2),
			fmt.Sprintf("%s=tryCatch(%s*%s,error=\\(e){NA})", newVar, var1, var2),
		}

		return newVar, ops[rand.Intn(len(ops))]
	}

	return "", "# No operation generated"
}
