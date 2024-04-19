"use client";

import Head from "next/head";
import Image from "next/image";
import { useEffect, useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { atomDark } from "react-syntax-highlighter/dist/cjs/styles/prism";
import { motion } from "framer-motion";
import { DiffEditor, Editor } from "@monaco-editor/react";
import EthIcon from "./ethereum-eth-logo.svg";

type OptimizationOptions = {
  structPacking: boolean;
  storageVariableCaching: boolean;
  callData: boolean;
};

// Option name
const optimizationOptionsNames: { [K in keyof OptimizationOptions]: string } = {
  structPacking: "Pack Structs",
  storageVariableCaching: "Cache Storage Variables",
  callData: "Optimise Call Data",
};

function getOptionName<K extends keyof OptimizationOptions>(option: K): string {
  return optimizationOptionsNames[option];
}

export default function Home() {
  const [inputCode, setInputCode] = useState("");
  const [unoptimizedCode, setUnoptimizedCode] = useState("");
  const [optimizedCode, setOptimizedCode] = useState("");
  const [testCode, setTestCode] = useState("");
  const [enableDiff, setEnableDiff] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [isErrorVisible, setIsErrorVisible] = useState(false);
  const [optimizationOptions, setOptimizationOptions] =
    useState<OptimizationOptions>({
      structPacking: false,
      storageVariableCaching: false,
      callData: false,
    });

  useEffect(() => {
    if (error) {
      setIsErrorVisible(true);
      const timer = setTimeout(() => setIsErrorVisible(false), 5000);
      return () => clearTimeout(timer); // Cleanup on unmount
    }
  }, [error, 5000]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");

    try {
      const response = await fetch("/api/optimize", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          contractCode: inputCode,
          testCode: testCode,
          opts: optimizationOptions,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setOptimizedCode(data.data.optimizedCode);
        setUnoptimizedCode(data.data.unoptimizedCode);
      } else {
        const errorData = await response.json();
        setError(errorData.error);
      }
    } catch (error: any) {
      console.error("Error:", error);
      setError("An unexpected error occurred.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleOptimizationOptionChange = (
    option: keyof OptimizationOptions,
  ) => {
    setOptimizationOptions((prevOptions) => ({
      ...prevOptions,
      [option]: !prevOptions[option],
    }));
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: { duration: 0.1, staggerChildren: 0.1 },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { opacity: 1, y: 0 },
  };

  const functionList: functionSignature[] = [
    {
      name: "SUM",
      args: ["arg1", "arg2"],
    },
    {
      name: "functionName2",
      args: ["arg1", "arg2"],
    },
  ];

  const functionComponent = functionCallItem(functionList[0]);
  const defaultWarning =
    "// Enter your Solidity code here\n// Make sure that the code is syntactically correct\n// This tool is still in development and may not work as expected\n// Please use at your own risk\n// If you encounter any issues, please report them on the GitHub repository\n//";

  return (
    <motion.div
      className="min-h-screen bg-gradient-to-br from-slate-800 to-gray-700 text-white text-sm"
      variants={containerVariants}
      initial="hidden"
      animate="visible"
    >
      {TopBar()}
      {isErrorVisible && (
        <motion.div
          className="bg-red-500 text-white px-4 py-2 text-sm"
          variants={itemVariants}
          initial={{ opacity: 0, y: -10 }} // Initial state (hidden)
          animate={{ opacity: 1, y: 0 }} // Animation on mount
          exit={{ opacity: 0, y: -10 }} // Animation on dismount
          transition={{ duration: 0.1 }} // Smooth animation
        >
          {error}
        </motion.div>
      )}
      <motion.form onSubmit={handleSubmit} variants={itemVariants}>
        <div className="flex space-x-2 px-2 pt-2">
          <div className="flex flex-col w-1/3 px-4 bg-stone-900 justify-between">
            <div className="space-y-2 text-sm">
              <label className="block font-bold py-2 border-b border-gray-700">
                Optimizations
              </label>
              {Object.entries(optimizationOptions).map(([option, enabled]) => (
                <motion.div
                  key={option}
                  variants={itemVariants}
                  className={`flex items-center space-x-2 ${enabled ? "text-green-400" : "text-white"} cursor-pointer transition duration-300 hover:text-green-400`}
                  onClick={() =>
                    handleOptimizationOptionChange(
                      option as keyof OptimizationOptions,
                    )
                  }
                >
                  <span>
                    {enabled ? "✅ " : "⚡️ "}
                    {getOptionName(option as keyof OptimizationOptions)}
                  </span>
                </motion.div>
              ))}
            </div>
            <div className="text-center my-2 text-sm">
              <motion.button
                type="submit"
                className="w-full bg-blue-600 text-white py-1 font-bold transition duration-300 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
                disabled={isLoading}
                whileTap={{ scale: 0.95 }}
              >
                {isLoading ? "Optimizing..." : "Optimize Code"}
              </motion.button>
            </div>
          </div>
          <div className="w-full flex-col bg-stone-900">
            <h3 className="px-4 border-b border-gray-700 font-bold py-2">
              Input
            </h3>
            <Editor
              height="60vh"
              defaultLanguage="sol"
              language="sol"
              theme="vs-dark"
              defaultValue={defaultWarning}
              onChange={(value) =>
                value != undefined
                  ? setInputCode(value)
                  : console.log("undefined")
              }
            />
          </div>
          {/* New form for function and arguments */}
          <div className="flex-col w-1/3 mt-8 mb-6 px-4 border border-gray-600 rounded-lg hidden">
            <div className="flex">
              <div className="w-1/2 mr-2">
                <label
                  htmlFor="functionName"
                  className="block mb-2 font-bold"
                ></label>
                <input
                  id="functionName"
                  type="text"
                  className="w-full h-3 p-4 bg-gray-700 text-white border border-gray-600 rounded-lg transition duration-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="functionName"
                />
              </div>

              <div className="w-1/2 hidden">
                <label
                  htmlFor="functionArgs"
                  className="block mb-2 font-bold"
                ></label>
                <input
                  id="functionArgs"
                  type="text"
                  className="w-full h-3 p-4 bg-gray-700 text-white border border-gray-600 rounded-lg transition duration-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="args"
                  // Add validation logic for comma-separated values and data types (optional)
                />
              </div>
            </div>
            <button
              type="button"
              className="ml-4 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              Estimate
            </button>
          </div>
        </div>
      </motion.form>

      {optimizedCode && (
        <button
          onClick={() => setEnableDiff(!enableDiff)}
          className="absolute right-0 bg-blue-600 text-sm text-white mt-2 mx-2 p-1 px-2 font-semibold duration-300 hover:bg-blue-700 focus:outline-none "
        >
          Toggle Diff
        </button>
      )}
      {optimizedCode && !enableDiff && (
        <motion.div
          className="mt-2"
          variants={containerVariants}
          initial="hidden"
          animate="visible"
        >
          <motion.div
            className="grid grid-cols-1 md:grid-cols-2 gap-2 px-2"
            variants={itemVariants}
          >
            <div>
              <h3 className="text-xl font-bold mb-2">Unoptimized Code</h3>
              <SyntaxHighlighter language="solidity" style={atomDark}>
                {unoptimizedCode}
              </SyntaxHighlighter>
            </div>
            <div>
              <h3 className="text-xl font-bold mb-2">Optimized Code</h3>
              <SyntaxHighlighter language="solidity" style={atomDark}>
                {optimizedCode}
              </SyntaxHighlighter>
            </div>
          </motion.div>
        </motion.div>
      )}
      {optimizedCode && enableDiff && (
        <motion.div
          className="p-2"
          variants={containerVariants}
          initial="hidden"
          animate="visible"
        >
          <motion.div variants={itemVariants}>
            <h3 className="text-xl font-bold mb-2">Diff View</h3>
            <DiffEditor
              height="60vh"
              language="sol"
              theme="vs-dark"
              original={unoptimizedCode}
              modified={optimizedCode}
            />
          </motion.div>
        </motion.div>
      )}
    </motion.div>
  );
}

interface functionSignature {
  name: string;
  // array of args
  args: string[];
}

function functionCallItem(fs: functionSignature) {
  return (
    <div>
      <h4>{fs.name}</h4>
      <ul>
        {fs.args.map((arg, index) => (
          <div key={index} className="flex items-center">
            <input
              id={`arg-${index}`}
              type="text"
              className="w-full p-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder={arg} // Use the argument value as the placeholder
            />
          </div>
        ))}
      </ul>
    </div>
  );
}

function TopBar() {
  return (
    <header className="border-b border-solid border-gray-700 text-white p-4 flex justify-between items-center">
      <div className="text-sm font-bold flex">
        <div className="mr-4 ">
          <Image src={EthIcon} alt="eth-icon" width={12} height={12} />
        </div>
        <div>Solidity Gas Optimizer</div>
      </div>
      <div className="flex space-x-4 text-sm">
        <a href="#" className="hover:text-gray-400">
          Optimizer
        </a>
        <a href="#" className="hover:text-gray-400">
          Github
        </a>
        <a href="#" className="hover:text-gray-400">
          Contact
        </a>
      </div>
    </header>
  );
}
