# Solidity-Gas-Optimizoor

![Alt text](demo.png)

An high performance automated tool that optimizes gas usage in Solidity smart contracts, focusing on storage and function call efficiency.

For more information on architecture and implementation, see the [docs](docs.md)

**Disclaimer:**
This code is a work in progress and can contain bugs. Use it at your own risk. Feature request and bug reports are welcome.

**Problem Statement:**

Solidity developers need tools to help them write gas-efficient code to minimize the execution cost of smart contracts on the blockchain. While there are some linters and optimizers available, there's a lack of tools specifically designed to analyze and suggest optimizations at both the source code and intermediate representation levels.

**Project Objective:**

The goal of this project is to design and implement a CLI and Web app that analyzes Solidity smart contracts, identifies patterns that lead to high gas usage, and suggests or automatically applies optimizations to improve gas efficiency.

## How does our tool optimize the gas of your smart contracts?

### Struct Data Packing

- **Concept**: Aligning struct members under 32 bytes together optimizes storage usage on the EVM.
- **Advantages**: This technique minimizes the number of `SLOAD` or `SSTORE` operations, slashing storage interaction costs by 50% or more when dealing with multiple struct values within a single slot.
- **Documentation**: [Structured Data Packing Guidance](https://github.com/beskay/gas-guide/blob/main/OPTIMIZATIONS.md#storage-packing)

[Demo](https://www.youtube.com/watch?v=cm8SGd2WK24&ab_channel=YongKangChia)

### Caching Storage Variables

- **Approach**: Utilize local variables to cache frequently accessed storage variables, reducing the number of expensive storage reads and writes.
- **Details**: Create a temporary local variable to store the value of a storage variable if it's accessed multiple times.
- **Source**: [Caching Storage Variables](https://www.rareskills.io/post/gas-optimization#viewer-8lubg)

### Calldata Efficiency

- **Gas Savings**: Leveraging calldata for unaltered external function inputs is more cost-effective than utilizing memory.
- **Implementation**: Analyze functions to ensure that inputs declared as `memory` are not modified. If unmodified, convert to `calldata`.
- **Reference**: [Calldata Efficiency Tips](https://github.com/beskay/gas-guide/blob/main/OPTIMIZATIONS.md#calldata-instead-of-memory-for-external-functions)

---

## Setup

### Golang

download dependencies: `go mod download`

update submodules: `make update-submodules`

### Web

1. Install npm dependencies: `npm install`
2. Install tmux with `sudo apt install tmux` or `brew install tmux`
3. Start the backend and frontend with `make start`

### Solc

To get Solidity compiler releases and ensure that the `{HOME}/.solc/releases` directory exists, you can follow these steps:

1. Open a terminal or command prompt.

2. Create the directory structure for Solidity compiler releases by running the following command:

   ```
   mkdir -p ~/.solc/releases
   ```

   This command will create the `releases` directory inside `~/.solc` (`{HOME}/.solc`) if it doesn't already exist.

3. Download the desired Solidity compiler release from the official Solidity releases page on GitHub:

   - Go to the Solidity releases page: [https://github.com/ethereum/solidity/releases](https://github.com/ethereum/solidity/releases)
   - Choose the release version you want to download (e.g., `v0.8.4`).
   - Scroll down to the "Assets" section and download the appropriate binary for your operating system (e.g., `solc-windows.exe`, `solc-macos`, `solc-static-linux`).

4. Move the downloaded Solidity compiler binary to the `{HOME}/.solc/releases` directory. For example, if you downloaded `solc-macos` for macOS, you can move it using the following command:

   ```
   mv ~/Downloads/solc-macos ~/.solc/releases/solc-v0.8.4
   ```

   Make sure to replace `v0.8.4` with the actual version you downloaded.

5. (Optional) Create a symbolic link to the downloaded Solidity compiler binary for easier access. For example:

   ```
   ln -s ~/.solc/releases/solc-v0.8.4 ~/.solc/solc
   ```

   This creates a symbolic link named `solc` in the `~/.solc` directory that points to the specific version of the Solidity compiler you downloaded.

6. Verify that the Solidity compiler is installed correctly by running the following command:
   ```
   ~/.solc/solc --version
   ```
   It should display the version of the Solidity compiler you downloaded.

By following these steps, you will have downloaded the desired Solidity compiler release and ensured that the `{HOME}/.solc/releases` directory exists with the compiler binary inside it.

Note: Make sure to replace `{HOME}` with the actual path to your home directory if necessary.

## Design

Overall Architecture:

![Architecture](image-2.png)

- Each component is designed as a rust library. The main program is in `cli` which will be the main entrypoint that performs these optimisations.

### Lexer

Functionality:

- Lexer is the component that takes raw input text and converts it into a stream of tokens. Tokens are the basic building blocks of a language's syntax, such as keywords, identifiers, literals, operators, and punctuation symbols.

To handle struct packing and calldata optimizations, Lexer should recognize and generate tokens for:

- Data type declarations (e.g., uint, address, struct)
- Storage qualifiers (e.g., memory, calldata, storage)
- Function definitions and parameters
- Variable declarations and assignments
- Comments and whitespace (to keep intact)

### Parser

The parser takes the stream of tokens and builds an Abstract Syntax Tree (AST), a tree-like representation of the syntactic structure of the code.

For the optimizations, the parser should:

1. Build nodes for struct definitions, capturing the order and types of fields.
2. Build nodes for function definitions, categorizing them by visibility (public, external, etc.) and whether they are read-only or state-changing.
3. Recognize repeated access to storage variables, to identify opportunities for caching.
   The parser must handle Solidity's grammar accurately to construct a correct AST, which will be traversed during the optimization phase.

### Optimizer

The optimizer traverses the AST and applies transformations to optimize the code.

**For Struct Packing:**

Logic for struct packing:

1. Get the fields from struct
2. Sort the fields in decreasing order of size
3. Pack the fields into storage slots

Pseudocode for field packing:

- Initialization: The function starts with a list of fields to pack and an initially empty or partially filled list of StorageSlot bins. It also initializes an empty list to hold different packing options (packing_options), representing various ways the fields can be packed into slots.

- The function takes the first field from the list (the one to pack next) and checks each existing slot to see if the field can fit.

- If the field fits within the slot (the combined size of the field plus the slot's current offset is less than or equal to 32 bytes), the function:
  - Creates a copy of the current list of slots.
  - Adds the field to the appropriate slot in the copied list.
  - Adjusts the offset of the slot to account for the added field's size.
  - Recursively calls bin_packing with the remainder of the fields and the updated list of slots, then adds the returned packing configuration to the packing_options.
- If the field does not fit in any existing slot, the function:

  - Creates a new StorageSlot bin and places the current field in it.
  - Recursively calls bin_packing (grouping algorithm) with the remainder of the fields and the updated list of slots, including the new slot, then adds the returned packing configuration to the packing_options.

- After processing all fields, the function returns the packing_options list, which contains all possible packing configurations.

- The main function then selects the best packing configuration from the list of options and applies it to the struct definition.

**For Storage Variable Caching:**

1. Identify functions with multiple reads to the same storage variable.
2. Introduce a local variable at the beginning of the function to cache the storage read.
3. Replace subsequent reads with references to the cached local variable.
4. Ensure the caching does not interfere with any writes to the storage variable within the function scope.

For Calldata Optimization:

1. Identify external functions with parameters declared as memory.
2. Analyze the function body to check if the memory parameters are modified.
3. If no modifications are detected, change the parameter type to calldata.

### Printer

After optimization, the transformed AST must be converted back into Solidity source code. This involves:

- Writing a code generator that traverses the optimized AST and produces Solidity code.
- Ensuring comments and formatting are preserved as much as possible.
- Verifying that the generated code compiles and behaves as intended.
