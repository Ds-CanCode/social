"use client";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Footer } from "./page";
import { WebSocketProvider } from "../utilis/websocket.js";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <head>
        <link rel="icon" href="./favicon.png" />
      </head>
      <body className={`${geistSans.variable} ${geistMono.variable}`}>
        <WebSocketProvider>
          {children}
        </WebSocketProvider>
      </body>
    </html>
  );
}