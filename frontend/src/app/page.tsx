"use client";

import Head from "next/head";
import { useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { atomDark } from "react-syntax-highlighter/dist/cjs/styles/prism";

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
      const errorMessage = error.message || "An unexpected error occurred.";
      setError(errorMessage);
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

  return (
    <div className="min-h-screen bg-gray-900 text-white flex items-center justify-center">
      <Head>
        <title>Solidity Code Optimizer</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <div className="max-w-4xl w-full bg-gray-800 p-8 rounded shadow-md">
        <h1 className="text-3xl font-bold mb-8 text-center">
          Solidity Code Optimizer
        </h1>

        {error && (
          <div className="bg-red-500 text-white px-4 py-2 mb-4 rounded">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-8">
            <label htmlFor="inputCode" className="block mb-2 font-bold">
              Enter Solidity Code:
            </label>
            <textarea
              id="inputCode"
              className="w-full h-64 p-4 bg-gray-700 text-white border border-gray-600 rounded"
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
                <div key={option}>
                  <label>
                    <input
                      type="checkbox"
                      checked={enabled}
                      onChange={() =>
                        handleOptimizationOptionChange(
                          option as keyof OptimizationOptions
                        )
                      }
                    />
                    <span className="ml-2">{option}</span>
                  </label>
                </div>
              ))}
            </div>
          </div>

          <div className="text-center">
            <button
              type="submit"
              className="bg-blue-600 text-white px-6 py-3 rounded font-bold text-lg"
              disabled={isLoading}
            >
              {isLoading ? "Optimizing..." : "Optimize Code"}
            </button>
          </div>
        </form>

        {optimizedCode && (
          <div className="mt-12">
            <h2 className="text-2xl font-bold mb-6 text-center">
              Optimized Code
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              <div>
                <h3 className="text-xl font-bold mb-4">Original Code</h3>
                <SyntaxHighlighter language="solidity" style={atomDark}>
                  {inputCode}
                </SyntaxHighlighter>
              </div>
              <div>
                <h3 className="text-xl font-bold mb-4">Optimized Code</h3>
                <SyntaxHighlighter language="solidity" style={atomDark}>
                  {optimizedCode}
                </SyntaxHighlighter>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
