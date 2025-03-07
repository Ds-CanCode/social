'use client'

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
    faUsers, faArrowRight, faUserGroup,
    faTimes,
    faUserPlus,
    faSearch,
    faNewspaper,
    faCalendarAlt,
    faThumbsUp,
    faThumbsDown,
} from '@fortawesome/free-solid-svg-icons';
import { CreateGroupPost } from "../page.js"
import { ChatApplication } from "@/utilis/component/ChatApplication";
import { Leftsidebar } from "@/utilis/component/leftsidebar";
import { Navbar } from "@/utilis/component/navbar";
import { useState } from 'react';
import { use } from "react";
import { useEffect } from "react";
import { fetchUserInfo } from "@/utilis/fetching_data";
import { Post } from "@/utilis/component/display_post";
import { notFound, useRouter } from "next/navigation";
import showPopupNotification from '@/utilis/component/notification.js';

export default function Group({ params }) {
    const [isMobileRightSidebarOpen, setIsMobileRightSidebarOpen] = useState(false);
    const [groupdata, setGroupdata] = useState([]);
    const router = useRouter();
    const [userdata, setUserdata] = useState(null);
    const [showInvitePopup, setShowInvitePopup] = useState(false);

    const [eventdata, setEventdata] = useState(null);
    const [activeTab, setActiveTab] = useState('posts');
    const [postdata, setPostdata] = useState(null);

    const resolvedParams = use(params);
    const groupId = resolvedParams.id;

    if (groupId == 0) {
        router.push("/notfound");
    }
    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetchUserInfo(`api/groups/get?groupId=${groupId}`);
                
                if (response.status === 400 || response.status === 404) {
                    router.push("/notfound");
                    
                }

                setGroupdata(response || []);

            } catch (error) {
                console.error("Error fetching data:", error);
            }
        };

        fetchData();
    }, [groupId]);

    useEffect(() => {
        async function getUserData() {
            const userdata = await fetchUserInfo("api/users/info");
            setUserdata(userdata); // Store the user data in state
        }
        getUserData();
    }, []);

    useEffect(() => {
        async function getPostData() {
            const postdata = await fetchUserInfo(`api/post/get?groupId=${groupId}`);
            setPostdata(postdata);
        }
        getPostData();
    }, [groupId]);

    useEffect(() => {
        async function getEventData() {
            const eventdata = await fetchUserInfo(`api/groups/getevents?groupId=${groupId}`);
            setEventdata(eventdata);
            console.log("Event data:", eventdata)
        }
        getEventData();
    }, [groupId]);

    
    const handleTabChange = (tab) => {
        setActiveTab(tab);
    };

    return (
        <div className="group-hero">
            <Navbar setIsMobileRightSidebarOpen={setIsMobileRightSidebarOpen} />
            <Leftsidebar />
            <ChatApplication />
            <div className="group-container">
                <div className="group-header">
                    <div className="group-cover"></div>
                    <div className="group-infos">
                        <div className="group-avatar">
                            {/* <FontAwesomeIcon icon={faUserGroup} size="2x" /> */}
                            <img
                                className='group-image'
                                src={groupdata?.image ? `/api/images?path=${groupdata?.image}` : "/default-img.jpg"}
                                alt='group-img'
                            />
                        </div>
                        <div className="group-details">
                            <div className="group-name-section">
                                <h1 className="group-name">{groupdata?.title}</h1>
                                <span className="group-badge">
                                    <FontAwesomeIcon icon={faUsers} size="sm" />
                                    <span>{groupdata?.nbr} members</span>
                                </span>
                            </div>
                            <p className="group-description">
                                {groupdata?.Descreption}
                            </p>
                            <div className="group-actions">

                                <Memberstatus status={groupdata?.memberstatus?.status} pop={() => setShowInvitePopup(true)} group={groupId} />

                            </div>
                        </div>
                    </div>
                </div>
                <CreateGroupPost userdata={userdata} status={groupdata?.memberstatus?.status} />
                <InvitePopup
                    groupId={groupId}
                    isOpen={showInvitePopup}
                    onClose={() => setShowInvitePopup(false)}
                />

                 {/* Tabs Navigation */}
                <div className="group-tabs">
                    <div 
                        className={`tab ${activeTab === 'posts' ? 'active' : ''}`} 
                        onClick={() => handleTabChange('posts')}
                    >
                        <FontAwesomeIcon icon={faNewspaper} />
                        <span>Posts</span>
                    </div>
                    <div 
                        className={`tab ${activeTab === 'events' ? 'active' : ''}`} 
                        onClick={() => handleTabChange('events')}
                    >
                        <FontAwesomeIcon icon={faCalendarAlt} />
                        <span>Events</span>
                    </div>
                </div>

                {/* Tab Content */}
                <div className="group-content-layout">
                    {/* Posts Tab Content */}
                    <div className={`tab-content ${activeTab === 'posts' ? 'active' : ''}`}>
                        {postdata && postdata.length > 0 ? (
                            postdata.map((post) => <Post key={post.id} post={post} />)
                        ) : (
                            <div className="empty-state">
                                <FontAwesomeIcon icon={faNewspaper} size="3x" />
                                <p>No posts available in this group yet</p>
                            </div>
                        )}
                    </div>

                    {/* Events Tab Content */}
                    <div className={`tab-content ${activeTab === 'events' ? 'active' : ''}`}>
                        {eventdata && eventdata.length > 0 ? (
                            eventdata.map((event) => <EventCard key={event.id} event={event} />)
                        ) : (
                            <div className="empty-state">
                                <FontAwesomeIcon icon={faCalendarAlt} size="3x" />
                                <p>No events scheduled in this group yet</p>
                                <button className="create-content-btn">Create an event</button>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>

    );
}

function Memberstatus({ status, pop, group }) {
    const handleJoinGrp = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch(`/api/groups/request?groupId=${group}`, {
                method: 'POST',
                credentials: "include",
            });
            if (!response.ok) {
                throw new Error(`Erreur: ${response.status}`);
            }

            showPopupNotification("Invitation envoyée");
        } catch {
            console.error("Error Join Group");
        }
    }

    const handleAccORejGrp = async (e, str) => {
        e.preventDefault();
        try {
            const response = await fetch(`/api/groups/${str}?groupId=${group}`, {
                method: 'POST',
                credentials: "include",
            });
            if (!response.ok) {
                throw new Error(`Erreur: ${response.status}`);
            }

            showPopupNotification("Invitation accepted");
        } catch {
            console.error("Error Join Group");
        }
    }

    switch (true) {
        case status === "pending" :
            return <>
            <button className="join-group"  onClick={(e) => handleAccORejGrp(e, "acceptinvitation")}>
               accepte
                <FontAwesomeIcon icon={faArrowRight} size="sm" />
            </button>
            <button className="join-group" onClick={(e) => handleAccORejGrp(e, "rejectrequest")} >
               reject
                <FontAwesomeIcon icon={faArrowRight} size="sm" />
            </button>
        </>
        case status === "joinRequest" :
            return <>
                <button className="join-group" >
                    Pending
                    <FontAwesomeIcon icon={faArrowRight} size="sm" />
                </button>
            </>
        case status === "can't find request":
            return <>
                <button className="join-group" onClick={handleJoinGrp}>
                    Join Group
                    <FontAwesomeIcon icon={faArrowRight} size="sm" />
                </button>
            </>
        case status === "creator" || status === "member":
            return <button onClick={pop} className='invite-members join-group'>invite members</button>


        default:
            return <button className="join-group" onClick={handleJoinGrp}>
                Join Group
                <FontAwesomeIcon icon={faArrowRight} size="sm" />
            </button>
    }

}

export function InvitePopup({ groupId, isOpen, onClose }) {
    const [users, setUsers] = useState([]);
    const [loading, setLoading] = useState(true);
    const [selectedUsers, setSelectedUsers] = useState([]);
    const [searchQuery, setSearchQuery] = useState('');
    const [inviteStatus, setInviteStatus] = useState({ show: false, message: '', isError: false });

    useEffect(() => {
        if (isOpen) {
            fetchUsers();
        }
    }, [isOpen]);

    const fetchUsers = async () => {
        setLoading(true);
        try {
            const data = await fetchUserInfo("api/users/followers"); // Use fetchUserInfo for consistency
            console.log(data)
            if (data && data.status !== 401) {
                // Map API response to match expected format
                const formattedUsers = data.map(user => ({
                    id: user.Id,
                    avatar: user.Avatar,
                    name: `${user.FirstName} ${user.LastName}`, // Prefer nickName, fallback to full name
                }));
                setUsers(formattedUsers);
            } else {
                console.error("Unauthorized or invalid response");
            }
        } catch (error) {
            console.error("Error fetching users:", error);
        } finally {
            setLoading(false);
        }
    };

    const handleUserSelect = (userId) => {
        setSelectedUsers(prev => {
            if (prev.includes(userId)) {
                return prev.filter(id => id !== userId);
            } else {
                return [...prev, userId];
            }
        });
    };

    const handleSearchChange = (e) => {
        setSearchQuery(e.target.value);
    };

    const filteredUsers = users.filter(user =>

        user.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        user.username?.toLowerCase().includes(searchQuery.toLowerCase())
    );

    const handleInvite = async (e) => {
        e.preventDefault();
        if (selectedUsers.length === 0) {
            setInviteStatus({
                show: true,
                message: 'Please select at least one user to invite',
                isError: true
            });
            return;
        }

        const data = new FormData();
        data.append("users", selectedUsers);
        data.append("group", groupId);

        try {
            const response = await fetch(`/api/groups/invitetogroup`, {
                method: 'POST',
                credentials: "include",
                body: data
            });

            if (response.ok) {
                setInviteStatus({
                    show: true,
                    message: 'Invitations sent successfully!',
                    isError: false
                });

                // Reset selections
                setSelectedUsers([]);

                // Auto close after success (optional)
                setTimeout(() => {
                    onClose();
                    setInviteStatus({ show: false, message: '', isError: false });
                }, 2000);
            } else {
                throw new Error('Failed to send invitations');
            }
        } catch (error) {
            console.error('Error sending invitations:', error);
            setInviteStatus({
                show: true,
                message: 'Failed to send invitations. Please try again.',
                isError: true
            });
        }
    };
    

    if (!isOpen) return null;

    return (
        <div className="invite-popup-overlay">
            <div className="invite-popup-content">
                <button className="invite-close-button" onClick={onClose}>
                    <FontAwesomeIcon icon={faTimes} />
                </button>

                <h2 className="invite-popup-title">
                    <FontAwesomeIcon icon={faUserPlus} />
                    Invite Members
                </h2>

                <div className="invite-search-container">
                    <FontAwesomeIcon icon={faSearch} className="search-icon" />
                    <input
                        type="text"
                        placeholder="Search users..."
                        value={searchQuery}
                        onChange={handleSearchChange}
                        className="invite-search-input"
                    />
                </div>

                {inviteStatus.show && (
                    <div className={`invite-status-message ${inviteStatus.isError ? 'error' : 'success'}`}>
                        {inviteStatus.message}
                    </div>
                )}

                <div className="invite-users-container">
                    {loading ? (
                        <div className="invite-loading">Loading users...</div>
                    ) : filteredUsers.length > 0 ? (
                        filteredUsers.map(user => (
                            <label key={user.id} className="invite-user-item">
                                <div className="invite-user-checkbox-container">
                                    <input
                                        type="checkbox"
                                        checked={selectedUsers.includes(user.id)}
                                        onChange={() => handleUserSelect(user.id)}
                                        className="invite-user-checkbox"
                                    />
                                    <span className="custom-checkbox"></span>
                                </div>

                                <div className="invite-user-avatar">
                                    <img
                                        src={user.avatar ? `/api/images?path=${user.avatar}` : "/default-user.jpg"}
                                        alt={user.name || user.username}
                                    />
                                </div>

                                <div className="invite-user-info">
                                    <span className="invite-user-name">{user.name || user.username}</span>
                                    {user.email && <span className="invite-user-email">{user.email}</span>}
                                </div>
                            </label>
                        ))
                    ) : (
                        <div className="invite-no-results">No users found matching "{searchQuery}"</div>
                    )}
                </div>

                <div className="invite-actions">
                    <span className="invite-selected-count">
                        {selectedUsers.length} {selectedUsers.length === 1 ? 'user' : 'users'} selected
                    </span>
                    <button
                        className="invite-submit-button"
                        onClick={handleInvite}
                        disabled={selectedUsers.length === 0}
                    >
                        Send Invitations
                        <FontAwesomeIcon icon={faArrowRight} />
                    </button>
                </div>
            </div>
        </div>
    );
}

const EventCard = ({ event }) => {
    // Format the date for display
    const formatDate = (dateString) => {
        const options = { 
            weekday: 'long', 
            year: 'numeric', 
            month: 'long', 
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        };
        return new Date(dateString).toLocaleString(undefined, options);
    };

    const groupId = window.location.pathname.split("/").pop();
    
    const [attendeesCount, setAttendeesCount] = useState(event?.countattends || 0);

    const handleJoinEvent = async () => {
        try {
            const response = await fetch(
                `/api/groups/handleJoinEvent?groupId=${groupId}&eventId=${event.id}`,
                {
                    method: "POST",
                    credentials: "include",
                }
            );
    
            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    setAttendeesCount(prev => prev + 1); // Increment count
                    showPopupNotification("You are now attending this event");
                    // alert("You are now attending this event");
                }
            } else {
                const errorData = await response.json();
                throw new Error(errorData.message || "Failed to join event");
            }
        } catch (error) {
            console.error("Error joining event:", error);
            showPopupNotification(error.message);
        }
    };
    
    const handleDeleteEvent = async () => {
        try {
            const response = await fetch(
                `/api/groups/handleDeleteEvent?groupId=${groupId}&eventId=${event.id}`,
                {
                    method: "POST",
                    credentials: "include",
                }
            );
    
            if (response.ok) {
                const data = await response.json();
                if (data.success) {
                    setAttendeesCount(prev => Math.max(0, prev - 1)); // Ensure count doesn't go below 0
                    showPopupNotification("Event deleted successfully");
                }
            } else {
                const errorData = await response.json();
                throw new Error(errorData.message || "Failed to delete event");
            }
        } catch (error) {
            console.error("Error deleting event:", error);
            showPopupNotification(error.message);
        }
    };
    

    return (
        <div className="event-card">
            <div className="event-header">
                <div className="event-creator">
                    <div className="event-creator-avatar">
                        <img 
                            src={event?.user?.path ? `/api/images?path=${event.user.path}` : "/default-avatar.jpg"} 
                            alt="Creator" 
                        />
                    </div>
                    <div className="event-creator-info">
                        <span className="event-creator-name">
                            {event?.user?.nickname || event?.user?.firstname || "Loading..."}
                        </span>
                        <span className="event-date">
                            {event?.created_at ? formatDate(event.created_at) : "Recent"}
                        </span>
                    </div>
                </div>
            </div>

            <div className="event-content">
                <h3 className="event-title">{event?.title || "Loading..."}</h3>
                <p className="event-description">{event?.Descreption || "Loading..."}</p>
                
                <div className="event-meta">
                    <div className="event-meta-item">
                        <FontAwesomeIcon icon={faCalendarAlt} />
                        <span>{event?.eventtime ? formatDate(event.eventtime) : "Date TBD"}</span>
                    </div>
                    
                    <div className="event-meta-item">
                        <FontAwesomeIcon icon={faUsers} />
                        {/* <span>{event?.attendeesCount || 0} attending</span> */}
                        <span>{attendeesCount} attending</span>

                    </div>
                </div>
            </div>
            
            <div className="event-actions">
                <button onClick={handleJoinEvent} className="event-action-btn">
                    <FontAwesomeIcon icon={faThumbsUp} />
                    <span>Attend</span>
                </button>
                <button onClick={handleDeleteEvent} className="event-action-btn">
                    <FontAwesomeIcon icon={faThumbsDown} />
                    <span>Not Interested</span>
                </button>
            </div>
        </div>
    )
}
