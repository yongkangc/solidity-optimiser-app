"use client";

import Head from "next/head";
import { useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { atomDark } from "react-syntax-highlighter/dist/cjs/styles/prism";

export default function Home() {
  const [inputCode, setInputCode] = useState("");
  const [optimizedCode, setOptimizedCode] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

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
        body: JSON.stringify({ contractCode: inputCode }),
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
    }

    setIsLoading(false);
  };

  const closeAlert = () => {
    setError("");
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
          <div className="bg-red-500 text-white px-4 py-2 mb-4 rounded flex items-center justify-between">
            <span>{error}</span>
            <button
              className="text-white hover:text-gray-200 ml-4"
              onClick={closeAlert}
            >
              &times;
            </button>
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
            <div className="grid grid-cols-2 gap-8">
              <div>
                <h3 className="text-xl font-bold mb-4 text-center">
                  Original Code
                </h3>
                <SyntaxHighlighter language="solidity" style={atomDark}>
                  {inputCode}
                </SyntaxHighlighter>
              </div>
              <div>
                <h3 className="text-xl font-bold mb-4 text-center">
                  Optimized Code
                </h3>
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
