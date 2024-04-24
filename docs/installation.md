# Installation

## Requirements

- `>= go1.22.2` [link](https://go.dev/doc/install)
- `>= node 18` [link](https://nodejs.org/en/download)
- `foundryup` [link](https://book.getfoundry.sh/getting-started/installation)
- `make`

## Clone the repository

```bash
git clone --recurse-submodules https://github.com/yongkangc/solidity-optimiser-app.git
```

## Install dependencies

**Using make**

```bash
make install-dependencies
```

**CLI**

```bash
# project root
go run main.go
```

**Frontend**

```bash
cd frontend
npm install
```

**Estimator**

```bash
  cd estimator
  foundryup
```

## Using a custom solc binary (Only for Apple Silicon)
To get Solidity compiler releases for the `estimator`, you can follow these steps:

1. Open a terminal or command prompt.

2. Download the desired Solidity compiler release from the official Solidity releases page on GitHub:

   - Go to the Solidity releases page: [https://github.com/ethereum/solidity/releases](https://github.com/ethereum/solidity/releases)
   - Choose the release version you want to download (e.g., `v0.8.4`).
   - Scroll down to the "Assets" section and download the appropriate binary for your operating system (e.g., `solc-windows.exe`, `solc-macos`, `solc-static-linux`).

4. Move the downloaded Solidity compiler binary to the `{PROJECT_ROOT}/estimator/` directory. For example, if you downloaded `solc-macos` for macOS, you can move it using the following command:

   ```
   mv ~/Downloads/solc-macos {PROJECT_ROOT}/estimator/solc
   ```

5. Edit `{PROJECT_ROOT}/estimator/foundry.toml` to use the solc binary:
  ```toml
  [profile.default]
  src = "src"
  out = "out"
  libs = ["lib"]
  # solc_version = "0.8.25"
  solc_version = "./solc"
  gas_reports = ["Unoptimized", "Optimized"]
  ```

By following these steps, you will have downloaded the desired Solidity compiler release and configured the `estimator` project to use it.
