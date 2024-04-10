import { NextResponse } from "next/server";

export async function POST(request: Request) {
  const { contractCode } = await request.json();

  // Call your backend API here to optimize the code
  const optimizedCode = await optimizeCode(contractCode);

  return NextResponse.json({ optimizedCode });
}

// Helper function to call your backend API
async function optimizeCode(contractCode: string) {
  const response = await fetch("http://localhost:8080/optimize", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ contractCode }),
  });

  if (response.ok) {
    const data = await response.json();
    return data.optimizedCode;
  } else {
    const error = await response.text();
    throw new Error("Optimization failed due to: " + error);
  }
}
