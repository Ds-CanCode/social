import { config } from "@fortawesome/fontawesome-svg-core";
import "@fortawesome/fontawesome-svg-core/styles.css"; // Prevents FOUC
config.autoAddCss = false; // Disable automatic CSS injection
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBars,
  faSearch,
} from "@fortawesome/free-solid-svg-icons";
import { useState, useEffect, useRef } from "react";
import { fetchUserInfo } from "../fetching_data";
import { useRouter } from "next/navigation";
import Link from "next/link.js";
import Image from "next/image.js";



//  ********* navbar *************//
export function Navbar({ setIsMobileRightSidebarOpen }) {
    const [isDropdownOpen, setIsDropdownOpen] = useState(false);
    const dropdownRef = useRef(null);
    const router = useRouter();
    const [userdata, setUserdata] = useState(null);
    useEffect(() => {
      async function getUserData() {
        const userdata = await fetchUserInfo("api/users/info");
        setUserdata(userdata); // Store the user data in state
      }
      getUserData();
    }, []);
  
    useEffect(() => {
      function handleClickOutside(event) {
        if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
          setIsDropdownOpen(false);
        }
      }
  
      document.addEventListener("mousedown", handleClickOutside);
      return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);
  
    const handleReport = () => {
      router.push("/report");
      setIsDropdownOpen(false);
    };
    const handleProfile = () => {
      router.push("/profile");
      setIsDropdownOpen(false);
    };
  
    const handle_search = (e) => {
      e.preventDefault();
      const searchInput = e.target.querySelector(".search-bar")
      if (searchInput.value === "") { 
        searchInput.style.border = "1px solid rgba(226, 24, 24, 0.8)";
        return;
      }    
      const searchType = document.querySelector('input[name="search-type"]:checked').value;
      const searchQuery = e.target.querySelector(".search-bar").value;
      router.push(`/search?query=${searchQuery}&type=${searchType}`);
    };
    const handleLogout = async () => {
      try {
        const response = await fetch("/api/logout", {
          method: "POST",
          credentials: "include",
        });
  
        if (!response.ok) {
          throw new Error("Logout failed");
        }
  
        if (response.status === 200) {
          console.log("Logout successful");
          router.push("/login");
        }
      } catch (error) {
        console.error("Error logging out:", error);
      } finally {
        setIsDropdownOpen(false);
      }
    };
  
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
            </Link>
            <div className="search-section">
              <div className="search-input">
                <FontAwesomeIcon icon={faSearch} className="search-icon" />
  
                <form onSubmit={handle_search} action="/search" method="get">
                  <input
                    type="text"
                    placeholder="Search..."
                    className="search-bar"
                  />
                </form>
  
                <div className="search-glow"></div>
              </div>
  
              <div className="radio-group">
                <label className="search-radio-label">
                  <input
                    type="radio"
                    name="search-type"
                    value="people"
                    className="search-radio-input"
                    defaultChecked
                  />
                  People
                </label>
                <label className="search-radio-label">
                  <input
                    type="radio"
                    name="search-type"
                    value="groups"
                    className="search-radio-input"
                  />
                  Groups
                </label>
              </div>
            </div>
          </div>
  
          <div className="profile-section" ref={dropdownRef}>
            <button
              className="mobile-sidebar-toggle"
              onClick={() => setIsMobileRightSidebarOpen((prev) => !prev)}
              aria-label="Toggle right sidebar"
            >
              <FontAwesomeIcon icon={faBars} />
            </button>
  
            <button
              className="profile-button"
              onClick={() => setIsDropdownOpen(!isDropdownOpen)}
            >
              <div className="profile-avatar">
                <img
                  src={
                    userdata?.avatar
                      ? `/api/images?path=${userdata.avatar}`
                      : "/default-img.jpg"
                  }
                  alt="Profile"
                  width={32}
                  height={32}
                  className="avatar-image"
                />
              </div>
              <span className="username">
                {userdata?.firstName || "loading..."}
              </span>
              <svg
                className={`dropdown-arrow ${isDropdownOpen ? "open" : ""}`}
                width="12"
                height="8"
                viewBox="0 0 12 8"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M1 1L6 6L11 1"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                />
              </svg>
            </button>
  
            {isDropdownOpen && (
              <div className="dropdown-menu">
                <button className="dropdown-item" onClick={handleProfile}>
                  Profile
                </button>
                <button className="dropdown-item" onClick={handleLogout}>
                  Logout
                </button>
                <button className="dropdown-item" onClick={handleReport}>
                  Report
                </button>
              </div>
            )}
          </div>
        </div>
      </nav>
    );
  }