'use client'

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
    faUsers, faArrowRight, faUserGroup,
    faTimes,
    faCalendarPlus,
    faPencilAlt,
    faImage
} from '@fortawesome/free-solid-svg-icons';
import "./group.css"
import { ChatApplication } from "@/utilis/component/ChatApplication";
import { Leftsidebar } from "@/utilis/component/leftsidebar";
import { Navbar } from "@/utilis/component/navbar";
import { useState, useEffect } from 'react';


export default function Group() {
    return (
        <div>
            <Navbar />
            <Leftsidebar />
            <ChatApplication />
            <div className='message_not_found'> this is not a specified group page!!</div>
        </div>
    )
}

export function CreateGroupPost({ userdata , status}) {
    if (status !== "creator" && status !== "member") {
        return <></>;
    }

    const [showOptionsPopup, setShowOptionsPopup] = useState(false);
    const [showEventForm, setShowEventForm] = useState(false);
    const [showPostForm, setShowPostForm] = useState(false);
    const [selectedImage, setSelectedImage] = useState(null);
    const [image, setImage] = useState(null); // For storing the image blob

    const [eventForm, setEventForm] = useState({
        title: "",
        description: "",
        date: "",
    });
    
    const [PostForm, setPostForm] = useState({
        title: "",
        content: "",
        privacy: "group",
    });

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setSelectedImage(URL.createObjectURL(file));
            setImage(file);
        }
    };

    const closeAllPopups = () => {
        setShowOptionsPopup(false);
        setShowEventForm(false);
        setShowPostForm(false);
    };

    const removeImage = () => {
        setSelectedImage(null);
        setImage(null);
        // Reset the file input
        const fileInput = document.getElementById("postImage");
        if (fileInput) fileInput.value = "";
    };

    const handlePostGrp = async (e) => {
        e.preventDefault();

        const formData = new FormData();
        formData.append("title", PostForm.title);
        formData.append("privacy", "group");
        formData.append("content", PostForm.content);
        const targetGroupId = window.location.pathname.split("/").pop();
        console.log(targetGroupId);
        
        formData.append("groupId", targetGroupId);

        if (image) formData.append("image", image);

        try {
            const response = await fetch("/api/post/add", {
                method: "POST",
                body: formData, // No need for Content-Type when using FormData
                credentials: "include",
            });

            if (response.status === 201) {
                console.log("Post created successfully");
                setShowPostForm(false);
                setPostForm({
                    title: "",
                    content: "",
                    privacy: "group", // Ensure it resets properly
                });
                removeImage();
            } else {
                console.log("Error response status:", response.status);
            }
        } catch (error) {
            console.error("Error during post creation:", error);
        }
    };

    const handleEventSubmit = async (e) => {
        e.preventDefault();

        const formData = new FormData();
        formData.append("title", eventForm.title);
        formData.append("description", eventForm.description);
        formData.append("date", eventForm.date);
        const targetGroupId = window.location.pathname.split("/").pop();
        formData.append("groupId", targetGroupId);
        console.log(formData);
        

        try {
            const response = await fetch("/api/groups/addevent", {
                method: "POST",
                body: formData, // No need for Content-Type when using FormData
                credentials: "include",
            });

            if (response.status === 201) {
                console.log("Event created successfully");
                setShowEventForm(false);
                setEventForm({
                    title: "",
                    description: "",
                    date: "",
                });
            } else {
                console.log("Error response status:", response.status);
            }
        } catch (error) {
            console.error("Error during event creation:", error);
        }
    };


    return (
        <div className='create-group-post'>
            <div className='post-input-container'>
                <div className='group-avatar01'>
                    <img
                        className='group-image'
                        src={userdata?.avatar ? `/api/images?path=${userdata?.avatar}` : "/default-img.jpg"}
                        alt='group-img'

                    />
                </div>
                <input
                    type="text"
                    placeholder="Start a discussion..."
                    onClick={() => setShowOptionsPopup(true)}
                    readOnly
                />
            </div>

            {/* Options Popup */}
            {showOptionsPopup && (
                <div className='popup-overlay' onClick={closeAllPopups}>
                    <div className='popup-content options-popup' onClick={e => e.stopPropagation()}>
                        <button className='close-button' onClick={closeAllPopups}>
                            <FontAwesomeIcon icon={faTimes} />
                        </button>
                        <div className='options-container'>
                            <button
                                className='option-button'
                                onClick={() => {
                                    setShowOptionsPopup(false);
                                    setShowEventForm(true);
                                }}
                            >
                                <FontAwesomeIcon icon={faCalendarPlus} />
                                Create a Group Event
                            </button>
                            <button
                                className='option-button'
                                onClick={() => {
                                    setShowOptionsPopup(false);
                                    setShowPostForm(true);
                                }}
                            >
                                <FontAwesomeIcon icon={faPencilAlt} />
                                Create a Group Post
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {/* Event Form Popup */}
            {showEventForm && (
                <div className="popup-overlay" onClick={closeAllPopups}>
                    <div className="popup-content form-popup" onClick={(e) => e.stopPropagation()}>
                        <button className="close-button" onClick={closeAllPopups}>
                            <FontAwesomeIcon icon={faTimes} />
                        </button>
                        <h2>Create Event</h2>
                        <form className="event-form" onSubmit={handleEventSubmit}>
                            <div className="form-group">
                                <label htmlFor="eventTitle">Event Title</label>
                                <input
                                    type="text"
                                    id="eventTitle"
                                    name="title"
                                    placeholder="Enter event title"
                                    value={eventForm.title} onChange={(e)=>setEventForm(prev => ({ ...prev, title: e.target.value }))}
                                />
                            </div>
                            <div className="form-group">
                                <label htmlFor="eventDescription">Event Description</label>
                                <textarea
                                    id="eventDescription"
                                    name="description"
                                    placeholder="Describe your event..."
                                    rows="4"
                                    value={eventForm.description} onChange={(e)=>setEventForm(prev => ({ ...prev, description: e.target.value }))}
                                ></textarea>
                            </div>
                            <div className="form-group">
                                <label htmlFor="eventDate">Event Date</label>
                                <input
                                    type="datetime-local"
                                    id="eventDate"
                                    name="date"
                                    value={eventForm.date} onChange={(e)=>setEventForm(prev => ({ ...prev, date: e.target.value }))}
                                />
                            </div>
                            <button type="submit" className="submit-button">
                                Create Event
                                <FontAwesomeIcon icon={faArrowRight} />
                            </button>
                        </form>
                    </div>
                </div>
            )}


            {/* Post Form Popup */}
            {showPostForm && (
                <div className='popup-overlay' onClick={closeAllPopups}>
                    <div className='popup-content form-popup' onClick={e => e.stopPropagation()}>
                        <button className='close-button' onClick={closeAllPopups}>
                            <FontAwesomeIcon icon={faTimes} />
                        </button>
                        <h2>Create Post</h2>
                        <form className='post-form' onSubmit={handlePostGrp}>
                            <div className='form-group'>
                                <label htmlFor="postTitle">Post Title</label>
                                <input type="text" id="postTitle" placeholder="Enter post title"  value={PostForm.title} onChange={(e)=>setPostForm(prev => ({ ...prev, title: e.target.value }))}/>
                            </div>
                            <div className='form-group'>
                                <label htmlFor="postImage">Image Upload</label>
                                <div className='image-upload-container'>
                                    <input
                                        type="file"
                                        id="postImage"
                                        accept="image/*"
                                        onChange={handleImageChange}
                                    />
                                    <label htmlFor="postImage" className='image-upload-label'>
                                        <FontAwesomeIcon icon={faImage} />
                                        <span>Choose Image</span>
                                    </label>
                                    {selectedImage && (
                                        <div className='image-preview'>
                                            <img src={selectedImage} alt="Preview" />
                                            <button
                                                type="button"
                                                className="remove-image"
                                                onClick={removeImage}
                                            >
                                                ×
                                            </button>
                                        </div>
                                    )}
                                </div>
                            </div>
                            <div className='form-group'>
                                <label htmlFor="postContent">Post Content</label>
                                <textarea
                                    id="postContent"
                                    placeholder="Write your post..."
                                    rows="4"
                                    value={PostForm.content}
                                    onChange={(e) => setPostForm(prev => ({ ...prev, content: e.target.value }))}
                                    ></textarea>
                            </div>
                            <button type="submit" className='submit-button'>
                                Create Post
                                <FontAwesomeIcon icon={faArrowRight} />
                            </button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}