import { NextResponse } from "next/server";

const API_AUTH_CHECK_URL = "http://localhost:8080/api/check-auth";
const AUTH_COOKIE_NAME = "session";

export async function middleware(request) {


  const url = request.nextUrl.clone();
  const fullPath = url.pathname + url.search; // Preserve query params

  console.log(`Rewriting: ${fullPath}`);

  if (url.pathname.startsWith('/api')) {
    // Rewrite request to the backend without modifying headers or body
    return NextResponse.rewrite(`http://localhost:8080${fullPath}`);
  }
  const authCookie = request.cookies.get(AUTH_COOKIE_NAME)?.value;
  const isAuthRoute = ["/login", "/register"].includes(request.nextUrl.pathname);
  
  if (!authCookie) {
    // If no session cookie and accessing protected routes -> Redirect to login
    if (!isAuthRoute) {
      return NextResponse.redirect(new URL("/login", request.url));
    }
    return NextResponse.next(); // Allow login/register pages
  }

  try {
    const response = await fetch(API_AUTH_CHECK_URL, {
      headers: {
        Cookie: `${AUTH_COOKIE_NAME}=${authCookie}`,
      },
      credentials: "include",
    });

    if (response.ok) {
      // If authenticated and accessing login/register -> Redirect to content
      if (isAuthRoute) {
        return NextResponse.redirect(new URL("/", request.url));
      }
    } else {
      // If authentication check fails -> Redirect to login
      if (!isAuthRoute) {
        return NextResponse.redirect(new URL("/login", request.url));
      }
    }
  } catch (error) {
    console.error("Auth check failed:", error);
    if (!isAuthRoute) {
      return NextResponse.redirect(new URL("/login", request.url));
    }
  }

  return NextResponse.next();
}
export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico|api/health|logo).*)"],
};


