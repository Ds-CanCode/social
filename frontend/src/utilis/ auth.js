// authHandler.js
export async function checkAuth() {
  try {
    const response = await fetch("/api/check-auth", {
      method: "GET",
      credentials: "include", // Include cookies in the request
    });

    if (response.status === 200) {
      const data = await response.json();
      return data.isAuthenticated; // Assuming the backend returns { isAuthenticated: true/false }
    } else {
      return false;
    }
  } catch (error) {
    console.error("Error checking authentication:", error);
    return false;
  }
}
