package optimizer

import "fmt"

func (o *Optimizer) optimizeStorageVariableCaching() {
	// TODO: check if the storage variable is read more than once
	//       if it is, cache it in a local variable
	contracts := o.builder.GetRoot().GetContracts()
	for _, contract := range contracts {
		// iterate through the contract's functions
		functions := contract.GetFunctions()
		for _, f := range functions {
			// iterate through the function's statements
			statements := f.GetAST().GetBody().GetStatements()
			for _, s := range statements {
				// check if the statement is a storage variable read
				fmt.Println(s)
			}
		}
	}

	// TODO: check if the storage variable is written more than once

	// TODO: Retrieve the type of the storage variable
	//       and check if it is a struct

	// TODO: Implement the caching of the storage variable

}
