import { Navbar } from "./navbar";
import { useState, useEffect, useRef } from "react";
import "./css/ChatApplication.css";
import { fetchUserInfo } from "../fetching_data";
import { useWebSocket } from "../websocket.js";
import { Send, X, Smile } from 'lucide-react';

const throttle = (func, delay) => {
  let lastCall = 0;
  return function (...args) {
    const now = new Date().getTime();
    if (now - lastCall < delay) {
      return;
    }
    lastCall = now;
    return func(...args);
  };
};

// ***************** this func is for rendering the right side-bar **************//
export function Rightsidebar({ isMobileOpen, onFriendClick, onGroupClick }) {
  const [friends, setFriends] = useState([]);
  const [groups, setGroups] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [friendsResponse, groupsResponse] = await Promise.all([
          fetchUserInfo(`api/users/userfollowers?profileId=0`),
          fetchUserInfo(`api/groups/getall`),
        ]);

        setFriends(friendsResponse);
        setGroups(groupsResponse);

        console.log("Friends:", friendsResponse);
        console.log("Groups:", groupsResponse);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    };

    fetchData();
  }, []);

  return (
    <aside className={`right-sidebar ${isMobileOpen ? "mobile-open" : ""}`}>
      <div className="right-sidebar-content">
        <div className="right-sidebar-section">
          <h3 className="section-title">Friends</h3>
          <div className="friends-container scrollable-container">
            {friends &&
              friends.map((friend) => (
                <div
                  key={friend.Id}
                  className="sidebar-item"
                  onClick={() => onFriendClick(friend.Id, friend.FirstName)}
                >
                  <div className="avatar-container">
                    <img
                      src={
                        friend?.Avatar
                          ? `/api/images?path=${friend.Avatar}`
                          : "/default-avatar.svg"
                      }
                      alt={`${friend.FirstName}'s avatar`}
                      className="avatar"
                      onError={(e) => {
                        e.target.onerror = null;
                        e.target.src = "/default-img.jpg";
                      }}
                    />
                    {/* <span className={`status-indicator ${friend.status}`}></span> */}
                  </div>
                  <span className="item-name">
                    {friend.FirstName + " " + friend.LastName}
                  </span>
                </div>
              ))}
          </div>
        </div>

        <div className="right-sidebar-section">
          <h3 className="section-title">Groups</h3>
          <div className="groups-container scrollable-container">
            {groups &&
              groups.map((group) => (
                <div
                  key={group.Id}
                  className="sidebar-item"
                  onClick={() => onGroupClick(group.Title, group.Id)}
                >
                  <div className="avatar-container">
                    <img
                      src={
                        group?.Path
                          ? `/api/images?path=${group.Path}`
                          : "/default-img.jpg"
                      }
                      alt={`${group.Title} group avatar`}
                      className="avatar group-avatar"
                    />
                  </div>
                  <div className="item-details">
                    <span className="item-name">{group.Title}</span>
                    <span className="item-meta">
                      {group.MemberCount} members
                    </span>
                  </div>
                </div>
              ))}
          </div>
        </div>
      </div>
    </aside>
  );
}

//********************** rendering the chatbox here it takes three params *************************** //
const EMOJIS = ["😊", "😂", "🥰", "🤓", "😎", "🤔", "👍", "❤️", "✨", "🎮", "🚀", "💻", "🌟", "🔥", "💪", "🎯"];

export function Chatbox({ activeChatuser, isVisible, setIsVisible }) {
  const [messages, setMessages] = useState([]);
  const [newMessage, setNewMessage] = useState("");
  const [showEmojiPicker, setShowEmojiPicker] = useState(false);
  const messagesEndRef = useRef(null);
  const messagesContainerRef = useRef(null);
  const socket = useWebSocket();
  const [offset, setOffset] = useState(0);
  const [isFetching, setIsFetching] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [isInitialLoad, setIsInitialLoad] = useState(true);
  const scrollPositionRef = useRef(null);
  
  const handleChat = (totas, id) => {
    if (!totas || totas.length === 0) {
      setHasMore(false);
      return;
    }
    console.log(totas);
    
    
    const newMessages = totas.map(data => {

      if (data.Type === "Private") {
        if (id === data.SenderID) {
          return {
            Content: data.Content,
            sender: "other",
            senderName: activeChatuser.name,
            timestamp: new Date(data.Created_at).toLocaleTimeString(),
          };
        } else {
          return {
            Content: data.Content,
            sender: "self",
            senderName: "You",
            timestamp: new Date(data.Created_at).toLocaleTimeString(),
          };
        }
      } else if (data.Type==="Group") {
        if (data.Vv == "You") {
          return {
            Content: data.Content,
            sender: "self",
            senderName: "You",
            timestamp: new Date(data.Created_at).toLocaleTimeString(),
          };

        } else {
          return {
          Content: data.Content,
          sender: "other",
          senderName: data.Sender.lastname,
          timestamp: new Date(data.Created_at).toLocaleTimeString(),

        }

      }
      }
    });

    // messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
    // messagesContainerRef.current.scrollIntoView({ behavior : "smooth"})

    setMessages(prevMessages => [...newMessages,...prevMessages]);
    
    // Update offset for pagination
    setOffset(prev => prev + totas.length);
    
    // If we received fewer messages than requested, we've reached the end
    if (totas.length < 10) {
      setHasMore(false);
    }
  };

  const fetchChatHistory = async (isInitial = false) => {
    if (isFetching || !hasMore) return;
    
    try {
      setIsFetching(true);
      console.log("Fetching messages with offset:", offset);
      
      // If we're loading chat history for a new user, reset states
      if (isInitial) {
        setOffset(0);
        setMessages([]);
        setHasMore(true);
        scrollPositionRef.current = messagesContainerRef.current.scrollHeight;

      }
      
      // Store current scroll position before loading more messages
      if (!isInitial && messagesContainerRef.current) {
        scrollPositionRef.current = messagesContainerRef.current.scrollHeight;
      }
      let response
      const currentOffset = isInitial ? 0 : offset;
      if (activeChatuser.type =="Group") {
        response = await fetch(
          `/api/chatgrouphistory?groupId=${activeChatuser.id}&offset=${currentOffset}`, 
          { credentials: "include" }
        );

      } else  if (activeChatuser.type=="Private") {

        response = await fetch(
         `/api/chathistory?recivierID=${activeChatuser.id}&offset=${currentOffset}`, 
         { credentials: "include" }
       );
      }



      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      data.reverse()
      if (data && data.length > 0) {
        const sortedData = [...data]
        handleChat(sortedData, activeChatuser.id);

      } else {
        setHasMore(false);
      }
    } catch (error) {
      console.error("Error loading messages:", error);
    } finally {
      setIsFetching(false);
      setIsInitialLoad(false);
    }
  };

  // Initial load when active chat user changes
  useEffect(() => {
    if (activeChatuser && activeChatuser.id) {
      setIsInitialLoad(true);
      setOffset(0);
      setMessages([]);
      setHasMore(true);
      fetchChatHistory(true);
    }
  }, [activeChatuser]);

  // Scroll to bottom on initial load
  useEffect(() => {
    if (isInitialLoad && messages.length> 0 ) {
      console.log("Scrolling to bottom...");
      console.log(messagesEndRef.current);
      
      setTimeout(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
      }, 100); // Small delay to allow rendering
    }
  }, [isInitialLoad , messages.length]);
  

  // Restore scroll position after loading older messages
  useEffect(() => {
    if (!isInitialLoad && scrollPositionRef.current && messagesContainerRef.current) {
      const newScrollHeight = messagesContainerRef.current.scrollHeight;
      const scrollDiff = newScrollHeight - scrollPositionRef.current;
      messagesContainerRef.current.scrollTop = scrollDiff;
      scrollPositionRef.current = null;
    }
  }, [messages, isInitialLoad]);

  // Handle real-time messages via WebSocket
  useEffect(() => {
    if (socket) {
      const handleMessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          let newMsg
          if (data.type === activeChatuser.type && data.type === "Private") {
            const message = data.data;

            if (activeChatuser.id === message.SenderID) {
               newMsg = {
                Content: message.Content,
                sender: "other",
                senderName: activeChatuser.name,
                timestamp: new Date().toLocaleTimeString(),
              };
              setMessages((prevMessages) => [...prevMessages, newMsg]);
              
              // Scroll to bottom for new incoming messages
              setTimeout(() => {
                messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
              }, 100);
            }
          } else if (data.type === activeChatuser.type && data.type === "Group") {
            const message = data.data;

             newMsg = {
              Content: message.Content,
              sender: "other",
              senderName: data.SenderName,
              timestamp: new Date().toLocaleTimeString(),
            };
            setMessages((prevMessages) => [...prevMessages, newMsg]);
            
            // Scroll to bottom for new incoming messages
            setTimeout(() => {
              messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
            }, 100);

          }
        } catch (error) {
          console.error("Error parsing WebSocket message:", error);
        }
      };

      socket.addEventListener("message", handleMessage);

      return () => {
        socket.removeEventListener("message", handleMessage);
      };
    }
  }, [socket, activeChatuser]);

  // Handle scrolling to load more messages
  const handleScroll = throttle(() => {
    if (!messagesContainerRef.current || isFetching || !hasMore) return;
    
    const { scrollTop } = messagesContainerRef.current;
    
    // Load more when scrolled near the top (within 50px)
    if (scrollTop < 300) {
      fetchChatHistory();
    }
  }, 300);

  // Setup scroll event listener
  useEffect(() => {
    const container = messagesContainerRef.current;
    if (container) {
      container.addEventListener('scroll', handleScroll);
      return () => {
        container.removeEventListener('scroll', handleScroll);
      };
    }
  }, [handleScroll, isFetching, hasMore]);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (newMessage.trim()) {

      let newMsg
      if (activeChatuser.type === "Group") {
        newMsg = {
          type: "Group",
          data: {
            Group : activeChatuser.id ,
            ReceiverID: activeChatuser.id,
            Content: newMessage,
          },
        };

      } else {
        newMsg = {
         type: "Private",
         data: {
           ReceiverID: activeChatuser.id,
           Content: newMessage,
         },
       };

      }


      const displayMsg = {
        Content: newMsg.data.Content,
        sender: "self",
        senderName: "You",
        timestamp: new Date().toLocaleTimeString(),
      };
      setMessages(prevMessages => [...prevMessages, displayMsg]);

      socket.send(JSON.stringify(newMsg));
      setNewMessage("");
      
      // Scroll to bottom after sending a message
      setTimeout(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
      }, 100);
    }
  };

  const handleEmojiClick = (emoji) => {
    setNewMessage(prev => prev + emoji);
    setShowEmojiPicker(false);
  };

  if (!isVisible) return null;

  return (
    <div className="chatbox">
      <div className="chat-header">
        <div className="chat-user-info">
          <div className="user-avatar">
          </div>
          <span className="user-name">{activeChatuser.name}</span>
          <button
            className="close-button"
            onClick={() => setIsVisible(false)}
            aria-label="Close chat"
          >
            <X size={20} />
          </button>
        </div>
      </div>

      <div className="messages-container" ref={messagesContainerRef}>
        {isFetching && (
          <div className="loading-indicator">Loading messages...</div>
        )}
        
        {messages && messages.map((message, index) => (
          <div
            key={index}
            className={`message-wrapper ${
              message.sender === "self" ? "message-self" : "message-other"
            }`}
          >
            <div className="message">
              <span className="sender-name">{message.senderName}</span>
              <p className="message-text">{message.Content}</p>
              <span className="message-timestamp">{message.timestamp}</span>
              <div className="message-glow"></div>
            </div>
          </div>
        ))}
        <div ref={messagesEndRef} />
      </div>

      <form className="chat-input-form" onSubmit={handleSubmit}>
        <button
          type="button"
          className="emoji-button"
          onClick={() => setShowEmojiPicker(!showEmojiPicker)}
          aria-label="Open emoji picker"
        >
          <Smile size={20} />
        </button>

        {showEmojiPicker && (
          <div className="emoji-picker">
            {EMOJIS.map((emoji, index) => (
              <button
                key={index}
                type="button"
                className="emoji-item"
                onClick={() => handleEmojiClick(emoji)}
              >
                {emoji}
              </button>
            ))}
          </div>
        )}
        <input
          type="text"
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Type your message..."
          className="chat-input"
        />
        <button
          type="submit"
          className="send-button"
          disabled={!newMessage.trim()}
        >
          <Send size={20} />
        </button>
      </form>
    </div>
  );
}

// **** this is the main function for both the right side-bar and the chatbox ****//
export function ChatApplication() {
  const [isChatVisible, setIsChatVisible] = useState(false);
  const [activeChatuser, setActiveChatuser] = useState({ id: 0, name: "" });
  const [isMobileRightSidebarOpen, setIsMobileRightSidebarOpen] =
    useState(false);

  const handleFriendClick = (ID, Name) => {
    setActiveChatuser({ type :"Private" , id: ID, name: Name });
    setIsChatVisible(true);
    // Close mobile sidebar after selection on mobile devices
    setIsMobileRightSidebarOpen(false);
  };

  const handleGroupClick = (groupName, groupId) => {
    setActiveChatuser({ type:"Group",id: groupId, name: groupName });
    setIsChatVisible(true);
    // Close mobile sidebar after selection on mobile devices
    setIsMobileRightSidebarOpen(false);
  };

  return (
    <div className="chat-application">
      <Navbar setIsMobileRightSidebarOpen={setIsMobileRightSidebarOpen} />

      <Rightsidebar
        isMobileOpen={isMobileRightSidebarOpen}
        onFriendClick={handleFriendClick}
        onGroupClick={handleGroupClick}
      />
      {isChatVisible ? <Chatbox
        activeChatuser={activeChatuser}
        isVisible={isChatVisible}
        setIsVisible={setIsChatVisible}
      /> : null}
      
    </div>
  );
}