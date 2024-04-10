"use client";

import Head from "next/head";
import { useState } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { atomDark } from "react-syntax-highlighter/dist/cjs/styles/prism";

export default function Home() {
  const [inputCode, setInputCode] = useState("");
  const [optimizedCode, setOptimizedCode] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const response = await fetch("/api/optimize", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ code: inputCode }),
      });

      if (response.ok) {
        const { optimizedCode } = await response.json();
        setOptimizedCode(optimizedCode);
      } else {
        console.error("Optimization failed");
      }
    } catch (error) {
      console.error("Error:", error);
    }

    setIsLoading(false);
  };

  return (
    <div className="container mx-auto p-4">
      <Head>
        <title>Solidity Code Optimizer</title>
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <h1 className="text-3xl font-bold mb-4">Solidity Code Optimizer</h1>

      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label htmlFor="inputCode" className="block mb-2">
            Enter Solidity Code:
          </label>
          <textarea
            id="inputCode"
            className="w-full h-40 p-2 border border-gray-300 rounded text-black"
            value={inputCode}
            onChange={(e) => setInputCode(e.target.value)}
          />
        </div>
        <button
          type="submit"
          className="bg-blue-500 text-white px-4 py-2 rounded"
          disabled={isLoading}
        >
          {isLoading ? "Optimizing..." : "Optimize Code"}
        </button>
      </form>

      {optimizedCode && (
        <div className="mt-8">
          <h2 className="text-2xl font-bold mb-4 text-black">
            Optimized Code:
          </h2>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <h3 className="text-xl font-bold mb-2 text-black">
                Original Code
              </h3>
              <SyntaxHighlighter language="solidity" style={atomDark}>
                {inputCode}
              </SyntaxHighlighter>
            </div>
            <div>
              <h3 className="text-xl font-bold mb-2 text-black">
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
  );
}
