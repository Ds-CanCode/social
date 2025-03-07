"use client"
import React from "react"
import "./profile.css"
import { useState, useEffect } from "react";

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faCog,
  faUserPlus,
  faTimes,
} from '@fortawesome/free-solid-svg-icons';
import { ChatApplication } from "@/utilis/component/ChatApplication";
import { Leftsidebar } from "@/utilis/component/leftsidebar";
import { Navbar } from "@/utilis/component/navbar";
import { fetchUserInfo } from "@/utilis/fetching_data.js";
import { Post } from "@/utilis/component/display_post";
import Link from "next/link";

export default function Profile() {
  const [isMobileRightSidebarOpen, setIsMobileRightSidebarOpen] = useState(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  // const [isFollowing, setIsFollowing] = useState(null);
  const [userdata, setUserdata] = useState(null);
  const [postdata, setPostdata] = useState(null);
  const [followers, setFollowers] = useState(null);
  const [following, setFollowing] = useState(null);

  const [followPopover, setFollowPopover] = useState(false);
  const [followingPopover, setFollowingPopover] = useState(false);

  const toggleSettingsPopup = () => setIsSettingsOpen(!isMobileRightSidebarOpen);

  useEffect(() => {

    async function fetchData() {
      try {
        const [userResponse, postResponse, userfollowers, userfollowing] = await Promise.all([
          fetchUserInfo(`api/users/info`),
          fetchUserInfo(`api/post/getuserpost?targetId=0`),
          fetchUserInfo(`api/users/userfollowers?profileId=0`),
          fetchUserInfo(`api/users/userfollowing?profileId=0`),
        ]);

        setUserdata(userResponse || []);
        setPostdata(postResponse || []);
        setFollowers(userfollowers || []);
        setFollowing(userfollowing || []);
        console.log("Following data:", following);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    }

    fetchData();
  }, []);

  const followersHandler = () => {
    setFollowPopover(!followPopover);
    if (followingPopover) setFollowingPopover(false); // Close following popover if open
  }

  const followingsHandler = () => {
    setFollowingPopover(!followingPopover);
    if (followPopover) setFollowPopover(false); // Close followers popover if open
  }

  // Close popovers when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      const followersPopover = document.getElementById('followers-popover');
      const followingPopover = document.getElementById('following-popover');

      if (followPopover && followersPopover && !followersPopover.contains(event.target) &&
        !event.target.closest('.stat-card[data-type="followers"]')) {
        setFollowPopover(false);
      }

      if (followingPopover && followingPopover && !followingPopover.contains(event.target) &&
        !event.target.closest('.stat-card[data-type="following"]')) {
        setFollowingPopover(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, [followPopover, followingPopover]);

  return (
    <div className="profile-hero">
      <Navbar setIsMobileRightSidebarOpen={setIsMobileRightSidebarOpen} />
      <Leftsidebar />
      <ChatApplication />
      <div className="profile-container">
        <div className="profile-header">
          <div className="profile-cover"></div>
          <div className="user_profile-info">
            <div className="user_profile-avatar">
              <img src={userdata?.avatar ? `/api/images?path=${userdata.avatar}` : "/default-img.jpg"} alt="Profile" />
            </div>
            <div className="profile-details">
              <div className="profile-name-section">
                <h1 className="profile-name">{userdata?.firstName + " " + userdata?.lastName || "loading..."}</h1>
                <span className="profile-badge">{userdata?.nickName || ""}</span>

              </div>
              <p className="profile-bio">{userdata?.aboutme || ""}</p>
              <p className="profile-bio">birthday: {userdata?.datebirth ? formatDate(userdata.datebirth) : "loading..."}</p>
            </div>
            <div className="profile-actions">
              {/* <button className="edit-profile">
                <FontAwesomeIcon icon={faUserPlus} size="sm" /> Follow
              </button> */}
              <button className="settings-button" onClick={toggleSettingsPopup}>
                <FontAwesomeIcon icon={faCog} size="sm" />
              </button>

              <SettingsPopup
                isOpen={isSettingsOpen}
                onClose={() => setIsSettingsOpen(false)}
              />
            </div>
          </div>
        </div>

        <div className="profile-stats">
          <div className="stat-card">
            <span className="stat-value">{userdata?.nbrPosts}</span>
            <span className="stat-label">posts</span>
          </div>
          <div className="stat-card" data-type="followers" onClick={followersHandler}>
            <span className="stat-value">{followers?.length || 0}</span>
            <span className="stat-label">followers</span>
          </div>
          <div className="stat-card" data-type="following" onClick={followingsHandler}>
            <span className="stat-value">{following?.length || 0}</span>
            <span className="stat-label">following</span>
          </div>
          <div className="stat-card">
            <span className="stat-value">{userdata?.type === true ? "Private" : "Public"}</span>
            <span className="stat-label">account</span>
          </div>
        </div>

        {/* *********** followers popover ********* */}
        {followPopover && (
          <div id="followers-popover" className="followers-popover">
            <div className="followers-header">
              <h3>Followers</h3>
              <button className="close-popover" onClick={followersHandler}>
                <FontAwesomeIcon icon={faTimes} />
              </button>
            </div>
            <div className="followers-list">
              {followers?.length > 0 ? (
                followers.map((follower, index) => (
                  <div key={index} className="follower-item">
                    <Link href={`/profile/${follower.Id}`}>
                      <div className="follower-info">
                        <span className="follower-name">{follower.FirstName}</span>
                        <span className="follower-username">{follower.LastName}</span>
                      </div>
                    </Link>
                  </div>
                ))
              ) : (
                <div className="no-followers">No followers yet</div>
              )}
            </div>
          </div>
        )}

        {/* **** following popover ********* */}
        {followingPopover && (
          <div id="followers-popover" className="followers-popover">
            <div className="followers-header">
              <h3>Following</h3>
              <button className="close-popover" onClick={followingsHandler}>
                <FontAwesomeIcon icon={faTimes} />
              </button>
            </div>
            <div className="followers-list">
              {following?.length > 0 ? (
                following.map((following, index) => (
                  <div key={index} className="follower-item">

                    <Link href={`/profile/${following.Id}`}>
                      <div className="follower-info">
                        <span className="follower-name">{following.FirstName}</span>
                        <span className="follower-username">{following.LastName}</span>
                      </div>
                    </Link>
                  </div>
                ))
              ) : (
                <div className="no-followers">No following yet</div>
              )}
            </div>
          </div>
        )}

        <div className="profile-content">
          <div className="content-section">
            <h2 className="section-title">Recent Activity</h2>
            <div className="created-Posts">
              {postdata && postdata.length > 0 ? (
                postdata.map((post) => <Post key={post.id} post={post} />)
              ) : (
                <p>No posts available</p> // Ou un spinner de chargement
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}


const SettingsPopup = ({ isOpen, onClose }) => {
  const [privacySetting, setPrivacySetting] = useState('private');

  if (!isOpen) return null;

  const handleOverlayClick = (e) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const formDataToSend = new FormData()
    formDataToSend.append("type", privacySetting);
    try {
      const response = await fetch(`/api/users/changeprofiletype`, {
        method: "POST",
        credentials: "include",
        body: formDataToSend
      });
      if (response.status === 200) {
        alert('Bien')
      }
    } catch {
      console.error("error")
    }
    onClose();
  };

  const handlePrivacyChange = (e) => {
    setPrivacySetting(e.target.value);
  };

  return (
    <div className="setting-popup-overlay" onClick={handleOverlayClick}>
      <div className="settings-popup">
        <div className="settings-content">
          <h2>Privacy Settings</h2>
          <form onSubmit={handleSubmit}>
            <div className="setting-radio-group">
              <label>
                <input
                  type="radio"
                  name="privacy"
                  value="private"
                  checked={privacySetting === 'private'}
                  onChange={handlePrivacyChange}
                />
                Private
              </label>
              <label>
                <input
                  type="radio"
                  name="privacy"
                  value="public"
                  checked={privacySetting === 'public'}
                  onChange={handlePrivacyChange}
                />
                Public
              </label>
            </div>
            <button type="submit" className="submit-settings">
              Submit
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};

function formatDate(isoDate) {
  const [year, month, day] = isoDate.split("T")[0].split("-");
  return `${day}/${month}/${year}`;
}