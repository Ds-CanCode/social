"use client";

import Link from "next/link.js";
import { useState, useEffect } from "react";
import "./css/notificatoin.css";
// import { useNavigate } from "react-router-dom";

export function Notification({ notificationRef, notificationActive }) {
    const [activeTab, setActiveTab] = useState("group-invitation");

    const handleTabClick = (tabId) => {
        setActiveTab(tabId);
    };

    return (
        <>
            <div
                ref={notificationRef}
                className={`notification-container ${notificationActive ? "active" : ""
                    }`}
            >
                <div className="notification-tabs">
                    <button
                        className={`notification-tab ${activeTab === "group-invitation" ? "active" : ""
                            }`}
                        onClick={() => handleTabClick("group-invitation")}
                    >
                        Groups Notifications
                    </button>

                    <button
                        className={`notification-tab ${activeTab === "follow-requests" ? "active" : ""
                            }`}
                        onClick={() => handleTabClick("follow-requests")}
                    >
                        Follow Requests
                    </button>

                    <button
                        className={`notification-tab ${activeTab === "group-events" ? "active" : ""
                            }`}
                        onClick={() => handleTabClick("group-events")}
                    >
                        Group Events
                    </button>
                </div>
                <div className="notification-content">
                    <GroupsNotif activeTab={activeTab} />
                    <FollowRequest activeTab={activeTab} />
                    <GroupsEvent activeTab={activeTab} />
                </div>
            </div>
        </>
    );
}



export function GroupsNotif({ activeTab }) {
    const [elements, setElements] = useState([]);

    const addElement = (data) => {
        if (!data) return

        setElements((prevElements) => {

            const newElements = data
                .filter((d) => d?.Sender?.Id && d?.Group?.Id) // Ensure valid structure
                .map((d, index) => {


                    if (d.Type === "groupInvitation") {
                        let grouplink = "/group/" + d.Group.Id;
                        return (
                            <div key={index} className="group-invitation">
                                <p className="group-invitation-message">
                                    {d.Group.Title} is inviting you! Please check it out.
                                </p>
                                <Link href={grouplink} className="group-invitation-link">
                                    View Invitation
                                </Link>
                            </div>
                        );
                    } else {
                        let name = `${d.Sender.FirstName} ${d.Sender.LastName}`;

                        let status;

                        if (d.Accepted === true) {
                            status = <p className="accepted-message">✅ Request Accepted</p>
                        } else {
                            status = <div className="notification-group-actions">
                                <button className="group-action-btn attend-btn"
                                    onClick={() => AcceptRejectRequest(1, d.Group.Id, d.Sender.Id)}>Accept</button>
                                <button className="group-action-btn ignore-btn"
                                    onClick={() => AcceptRejectRequest(2, d.Group.Id, d.Sender.Id)}>Ignore</button>
                            </div>
                        }
                        return (
                            <div key={index} className="group-invitation">
                                <p className="group-invitation-message">
                                    {name} wants to join your group {d.Group.Title}.
                                </p>
                                {status}
                            </div>
                        );


                    }
                })
                .filter(Boolean); // Remove null values

            return [...prevElements, ...newElements];
        });
    };

    useEffect(() => {
        if (activeTab === "group-invitation") {
            async function FetchData() {
                try {
                    const result = (await FetchNotif("invitations")) || [];
                    const result1 = (await FetchNotif("requestgroup")) || [];
                    const combinedResults = [...result, ...result1];

                    addElement(combinedResults);
                } catch (error) {
                    console.error("Error fetching notifications:", error);
                }
            }

            FetchData();
        } else {
            setElements([]); // Clear notifications when switching tabs
        }
    }, [activeTab]);

    return (
        <div className={`tab-panel ${activeTab === "group-invitation" ? "active" : ""}`} id="group-invitation">
            <div className="notification-divider">{elements}</div>
        </div>
    );
}

async function AcceptRejectRequest(type, groupId, target) {
    let url = ""
    if (type == 1) {
        url = `/api/groups/acceptrequest?groupId=${groupId}&target=${target}`;
    } else if (type == 2) {
        url = `/api/groups/rejectrequest?groupId=${groupId}&target=${target}`;
    }

    try {
        let response = await fetch(url, { credentials: "include" });
        console.log(response.status);
        

        if (response.status === 401) {
            // useNavigate("/register"); // Uncomment if using React Router
            return;
        }

        if (response.status !== 200) {
            console.warn("Request failed with status:", response.status);
            return;
        }
    } catch (error) {
        console.error("Fetch error:", error);
    }
}


export function FollowRequest({ activeTab }) {
    const [elements, setElements] = useState([]);

    const addElement = (data) => {
        if (!data) return
        const elementsArray = data.map((d, index) => {

            let imgPath = "/api/images?path=" + d.Sender.Path;

            let profile = "/profile/" + d.Sender.Id;
            let name = d.Sender.FirstName + " " + d.Sender.LastName;

            return (
                <div key={index} className="follow-request">
                    <div className="follow-avatar">
                        <div className="follow-avatar-placeholder">
                            <img
                                className="follow-avatar"
                                src={imgPath}
                                alt="User Avatar"
                            />
                        </div>
                    </div>
                    <div className="follow-info">
                        <div className="follow-name">{name}</div>
                        <div className="follow-username">{d.Sender.Nickname}</div>
                    </div>
                    <Link href={profile} className="visit-account">
                        Visit
                    </Link>
                </div>
            );
        });
        setElements(elementsArray); // Set the array of JSX elements to the state
    };

    useEffect(() => {
        if (activeTab === "follow-requests") {
            async function FetchData() {
                try {
                    const result = await FetchNotif("requestuser");
                    addElement(result);
                } catch (error) {
                    console.error("Error fetching follow requests:", error);
                }
            }
            FetchData();
        } else {
            setElements([]); // Reset the elements array when the tab changes
        }
    }, [activeTab]);

    return (
        <div
            className={`tab-panel ${activeTab === "follow-requests" ? "active" : ""}`}
        >
            {elements}{" "}
        </div>
    );
}

export function GroupsEvent({ activeTab }) {
    const [elements, setElements] = useState([]);

    const addElement = (data) => {
        if (!data) return


        const elementsArray = data.map((d, index) => {

            let imgPath = "/api/images?path=" + d.Sender.Path;


            let profile = "/profile/" + d.Sender.Id;
            let name = d.Sender.FirstName + " " + d.Sender.LastName;

            return (
                <div key={index} className="follow-request">
                    <div className="follow-avatar">
                        <div className="follow-avatar-placeholder">
                            <img
                                className="follow-avatar"
                                src={imgPath}
                                alt="User Avatar"
                            />
                        </div>
                    </div>
                    <div className="follow-info">
                        <div className="follow-name">{name}</div>
                        <div className="follow-username">{d.Sender.Nickname}</div>
                    </div>
                    <Link href={profile} className="visit-account">
                        Visit
                    </Link>
                </div>
            );
        });
        setElements(elementsArray); // Set the array of JSX elements to the state
    };

    useEffect(() => {
        if (activeTab === "group-events") {
            async function FetchData() {
                try {
                    const result = await FetchNotif("event");
                    addElement(result);
                } catch (error) {
                    console.error("Error fetching follow requests:", error);
                }
            }
            FetchData();
        } else {
            setElements([]);
        }
    }, [activeTab]);

    return (
        <div
            className={`tab-panel ${activeTab === "group-events" ? "active" : ""}`}
            id="group-events"
        >
            {elements}{" "}
        </div>
    );
}

async function FetchNotif(type) {
    let url = "/api/notif/" + type;


    try {
        let response = await fetch(url, {
            credentials: "include",
        });

        if (response.status === 401) {
            //   useNavigate("/register");
            return undefined;
        }

        if (response.status !== 200) {
            return undefined;
        }

        return await response.json();
    } catch (error) {
        console.error("Fetch error:", error);
        return undefined;
    }
}

export default function showPopupNotification(message, duration = 3000) {
    const notification = document.createElement('div');
    notification.className = 'popup-Notif';
    const messageParagraph = document.createElement('p');
    messageParagraph.textContent = message;
    notification.appendChild(messageParagraph);
    document.body.appendChild(notification);
    notification.offsetHeight;
    setTimeout(() => {
      notification.classList.add('closing');
      notification.addEventListener('animationend', () => {
        document.body.removeChild(notification);
      });
    }, duration);
    return notification;
  }

