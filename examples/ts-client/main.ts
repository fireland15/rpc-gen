import {
  login,
  getTodos,
  changePassword,
  getTodo,
} from "./client"; // NOTE: .js extension is required in ESM TS

async function main() {
  const ctx = {
    baseUrl: "http://localhost:5175",
    // credentials: "include", // uncomment if your API uses cookies
  };

  try {
    console.log("→ Logging in...");
    await login("hunter2", "alice", ctx);
    console.log("✓ Login successful");

    console.log("→ Fetching todos...");
    const todos = await getTodos(ctx);
    console.log("✓ Todos:", todos);

    console.log("→ Changing password...");
    const ok = await changePassword("newpass123", "hunter2", ctx);
    console.log("✓ Password changed:", ok);

    console.log("→ Fetching todo...");
    const todo = await getTodo(123, ctx);
    console.log("✓ Todo:", todo);

  } catch (err) {
    console.error("✗ Error occurred");

    if (err instanceof Error) {
      console.error(err.message);
    } else {
      console.error(err);
    }

    process.exitCode = 1;
  }
}

main();
