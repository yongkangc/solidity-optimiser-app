import { NextResponse } from "next/server";

export async function POST(request: Request) {
  const { testCode } = await request.json();

  // Call your backend API here to optimize the code
  const data = await estimateGas(testCode);

  return NextResponse.json({ data });
}

// Helper function to call your backend API
async function estimateGas(testCode: string) {
  const response = await fetch("http://localhost:8080/estimate", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ testCode }),
  });

  if (response.ok) {
    const data = await response.json();
    return data;
  } else {
    const error = await response.text();
    throw new Error("Optimization failed due to: " + error);
  }
}
