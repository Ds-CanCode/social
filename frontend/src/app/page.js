"use client";
import {PostContainer, PostList} from "@/utilis/component/display_post.js";
import { Leftsidebar } from "@/utilis/component/leftsidebar";
import { ChatApplication } from "@/utilis/component/ChatApplication";
import React, { use } from "react";
import { config } from "@fortawesome/fontawesome-svg-core";
import "@fortawesome/fontawesome-svg-core/styles.css"; // Prevents FOUC
config.autoAddCss = false; // Disable automatic CSS injection
import { useState } from "react";
import { Navbar } from "@/utilis/component/navbar";
import Link from "next/link"
import Image from "next/image";

export default function Home() {
   const [isMobileRightSidebarOpen, setIsMobileRightSidebarOpen] =
     useState(false);
   return (
     <div className="hero">
       <Navbar setIsMobileRightSidebarOpen={setIsMobileRightSidebarOpen} />
       <Leftsidebar />
       <ChatApplication/>
       <PostContainer />
       <div className="main-content">
         <PostList />
       </div>
     </div>
   );
}

export function Navbarrend({ NavButton }) {
  return (
    <nav className="navbar">
      <div className="navbar-container">
        <div className="logo-section">
          <Link href="/" className="logo-link">
            <Image 
              src="/logo.svg"
              alt="Logo"
              width={40}
              height={40}
              className="logo-image"
            />
            <span className="logo-text">UNION</span>
          </Link>
        </div>
        {NavButton && <NavButton />} {/* Render the passed component */}
      </div>
    </nav>
  );
}

export function Footer() { 
  return (
    <footer className="footer">
      <p>&copy; 2025 Union. All rights reserved.</p>
    </footer>
  )
}

