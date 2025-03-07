
export async function fetchUserInfo(path) {
  try {
    const response = await fetch(`/${path}`, {
      method: "GET",
      credentials: "include",
    });

    const status = response.status; 

    if (status === 200) {
      const data = await response.json();
      return data 
    } else {
      return { status, data: null }; 
    }
  } catch (error) {
    console.error("Error fetching user info:", error);
    return false;
  }
}
