import { useState, useEffect, useRef } from 'react';
import Link from "next/link.js";
import { config } from "@fortawesome/fontawesome-svg-core";
import "@fortawesome/fontawesome-svg-core/styles.css"; // Prevents FOUC
config.autoAddCss = false;
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faHouse,
  faBell,
  faUserFriends,
} from "@fortawesome/free-solid-svg-icons";
import { Notification } from './notification';
import "./css/leftsidebar.css"


//****** leftsidebar that containes the navigation and it 
//* displays in both sides leftside for large devices and in the bottom for small devices  *******/
export let notifBTN = undefined
export function Leftsidebar() {
  const [isCollapsed, setIsCollapsed] = useState(true);
  const [notificationActive, setNotificationActive] = useState(false);

  const notificationRef = useRef(null);
  const notificationBtnRef = useRef(null);

  const menuItems = [
    { icon: faHouse, id: "home-btn", label: "Home", href: "/" },
    { icon: faBell, id: "notification-btn", label: "Notification", href: "#" },
    { icon: faUserFriends, id: "group-btn", label: "Groups", href: "/groups" },
  ];

  // Handle click outside to close notification
  useEffect(() => {
    function handleClickOutside(event) {
      if (
        notificationRef.current &&
        !notificationRef.current.contains(event.target) &&
        notificationBtnRef.current &&
        !notificationBtnRef.current.contains(event.target)
      ) {
        setNotificationActive(false);
      }
    }

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);



  notifBTN = notificationBtnRef.current;
  const handleNotificationClick = (e) => {
    e.preventDefault();

    if (notifBTN) {
      notifBTN.style.color = "white"
    }
    setNotificationActive(!notificationActive);
  };

  return (
    <>
      <aside className={`sidebar ${isCollapsed ? "collapsed" : ""}`}>
        <div className="sidebar-content">
          <button
            className="collapse-button"
            onClick={() => setIsCollapsed(!isCollapsed)}
            aria-label={isCollapsed ? "Expand sidebar" : "Collapse sidebar"}
          >
            <svg
              width="20"
              height="20"
              viewBox="0 0 20 20"
              fill="none"
              className={`collapse-icon ${isCollapsed ? "rotated" : ""}`}
            >
              <path
                d="M15 7L10 12L5 7"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>

          <nav className="sidebar-nav">
            {menuItems.map((item, index) => (
              item.id === "notification-btn" ? (
                <a
                  key={index}
                  href="#"
                  className="nav-item"
                  onClick={handleNotificationClick}
                  ref={notificationBtnRef}
                >
                  <FontAwesomeIcon icon={item.icon} className="nav-icon" />
                  <span id={item.id} className="nav-label">{item.label}</span>
                  <div className="nav-glow"></div>
                </a>
              ) : (
                <Link key={index} href={item.href} className="nav-item">
                  <FontAwesomeIcon icon={item.icon} className="nav-icon" />
                  <span id={item.id} className="nav-label">{item.label}</span>
                  <div className="nav-glow"></div>
                </Link>
              )
            ))}
          </nav>
        </div>
      </aside>

      {/* Notification Container */}
      <Notification notificationRef={notificationRef} notificationActive={notificationActive} />
    </>
  );
}