"use client";
import { usePathname } from "next/navigation";
import { createContext, useContext, useEffect, useState } from "react";

// Create a context for the WebSocket
const WebSocketContext = createContext(null);

// Create a provider component
export function WebSocketProvider({ children }) {
  const path = usePathname();
  const [socket, setSocket] = useState(null);

  useEffect(() => {
    if (socket) {
      if (path === "/login" || path === "/register") {
        socket.close();
      }
    } else {
      if (path === "/login" || path === "/register") {
        return;
      }

      // Initialize WebSocket connection
      const newSocket = new WebSocket("/api/ws");

      newSocket.addEventListener("open", () => {
        console.log("WebSocket connected");
        setSocket(newSocket)
      });

      newSocket.addEventListener("close", () => {
        console.log("WebSocket disconnected");
        setSocket(null);
      })

      // newSocket.addEventListener("message", (event) =>
      //   handleNotif(event)
      // );
      newSocket.addEventListener("message", (event) => {
        handleNotif(event);
        const data = JSON.parse(event.data);

        if (data.Type !== "follow" && data.Type !== "groupRequest" && data.Type !== "groupInvitation") {
          showPopupNotification("you have a new message");
        }
      });
      
      // Cleanup on unmount
    }
  }, [path]);

  return (
    <WebSocketContext.Provider value={socket}>
      {children}
    </WebSocketContext.Provider>
  );
}

// Custom hook to access the WebSocket context
export function useWebSocket() {
  return useContext(WebSocketContext);
}
import { notifBTN } from "./component/leftsidebar";
import showPopupNotification from '@/utilis/component/notification.js';
function handleNotif(event) {
  console.log(event);

  try {
    const message = JSON.parse(event.data);
    console.log(message);

    if (message.Type === "follow") {
      notifBTN.style.color = "red"
      showPopupNotification("You have new follower")
    } else if (message.Type == "groupRequest") {
      notifBTN.style.color = "red"
      showPopupNotification("You have new request group")
    } else if (message.Type == "groupInvitation") {
      showPopupNotification("You have new invitation")
      notifBTN.style.color = "red"
    }
  } catch {
    console.log("Error parsing WebSocket message.");
  }
}
