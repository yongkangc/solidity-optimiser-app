import { NextResponse } from "next/server";

export async function POST(request: Request) {
  const { contractCode, opts } = await request.json();

  // Call your backend API here to optimize the code
  const data = await optimizeCode(contractCode, opts);

  return NextResponse.json({ data });
}

// Helper function to call your backend API
async function optimizeCode(contractCode: string, opts: OptimizationConfig) {
  const response = await fetch("http://localhost:8080/optimize", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ contractCode, opts }),
  });

  if (response.ok) {
    const data = await response.json();
    return data;
  } else {
    const error = await response.text();
    throw new Error("Optimization failed due to: " + error);
  }
}

type OptimizationConfig = {
  structPacking: boolean;
  storageVariableCaching: boolean;
  callData: boolean;
};
