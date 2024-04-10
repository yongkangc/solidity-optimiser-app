"use client";

import Head from "next/head";
import { useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { atomDark } from "react-syntax-highlighter/dist/cjs/styles/prism";
import { motion } from "framer-motion";

type OptimizationOptions = {
  structPacking: boolean;
  storageVariableCaching: boolean;
  callData: boolean;
};

export default function Home() {
  const [inputCode, setInputCode] = useState("");
  const [optimizedCode, setOptimizedCode] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [optimizationOptions, setOptimizationOptions] =
    useState<OptimizationOptions>({
      structPacking: false,
      storageVariableCaching: false,
      callData: false,
    });

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
          opts: optimizationOptions,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        setOptimizedCode(data.optimizedCode);
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
    option: keyof OptimizationOptions
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
      transition: { duration: 0.6, staggerChildren: 0.2 },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { opacity: 1, y: 0 },
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 to-blue-900 text-white flex items-center justify-center">
      <Head>
        <title>Solidity Code Optimizer</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <motion.div
        className="max-w-4xl w-full bg-gray-800 p-8 rounded-lg shadow-lg"
        variants={containerVariants}
        initial="hidden"
        animate="visible"
      >
        <motion.h1
          className="text-4xl font-bold mb-8 text-center"
          variants={itemVariants}
        >
          Solidity Code Optimizer
        </motion.h1>

        {error && (
          <motion.div
            className="bg-red-500 text-white px-4 py-2 mb-4 rounded"
            variants={itemVariants}
          >
            {error}
          </motion.div>
        )}

        <motion.form onSubmit={handleSubmit} variants={itemVariants}>
          <div className="mb-8">
            <label htmlFor="inputCode" className="block mb-2 font-bold">
              Enter Solidity Code:
            </label>
            <textarea
              id="inputCode"
              className="w-full h-64 p-4 bg-gray-700 text-white border border-gray-600 rounded-lg transition duration-300 focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={inputCode}
              onChange={(e) => setInputCode(e.target.value)}
            />
          </div>

          <div className="mb-8">
            <label className="block mb-2 font-bold">
              Optimization Options:
            </label>
            <div className="space-y-2">
              {Object.entries(optimizationOptions).map(([option, enabled]) => (
                <motion.div key={option} variants={itemVariants}>
                  <label className="flex items-center">
                    <input
                      type="checkbox"
                      checked={enabled}
                      onChange={() =>
                        handleOptimizationOptionChange(
                          option as keyof OptimizationOptions
                        )
                      }
                      className="form-checkbox h-5 w-5 text-blue-500 transition duration-300"
                    />
                    <span className="ml-2">{option}</span>
                  </label>
                </motion.div>
              ))}
            </div>
          </div>

          <div className="text-center">
            <motion.button
              type="submit"
              className="bg-blue-600 text-white px-8 py-3 rounded-lg font-bold text-lg transition duration-300 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
              disabled={isLoading}
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              {isLoading ? "Optimizing..." : "Optimize Code"}
            </motion.button>
          </div>
        </motion.form>

        {optimizedCode && (
          <motion.div
            className="mt-12"
            variants={containerVariants}
            initial="hidden"
            animate="visible"
          >
            <motion.h2
              className="text-3xl font-bold mb-6 text-center"
              variants={itemVariants}
            >
              Optimized Code
            </motion.h2>
            <motion.div
              className="grid grid-cols-1 md:grid-cols-2 gap-8"
              variants={itemVariants}
            >
              <div>
                <h3 className="text-2xl font-bold mb-4">Original Code</h3>
                <SyntaxHighlighter language="solidity" style={atomDark}>
                  {inputCode}
                </SyntaxHighlighter>
              </div>
              <div>
                <h3 className="text-2xl font-bold mb-4">Optimized Code</h3>
                <SyntaxHighlighter language="solidity" style={atomDark}>
                  {optimizedCode}
                </SyntaxHighlighter>
              </div>
            </motion.div>
          </motion.div>
        )}
      </motion.div>
    </div>
  );
}
