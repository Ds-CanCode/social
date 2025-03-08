"use client"
import React from "react"
import "./groups.css"
import { useState, useEffect } from "react";
import { fetchUserInfo } from "@/utilis/fetching_data.js"; 
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { 
    faUserGroup, 
    faTimes, 
    faImage,
    faArrowRight,
    faPlus
} from '@fortawesome/free-solid-svg-icons';
import { ChatApplication } from "@/utilis/component/ChatApplication";
import { Leftsidebar } from "@/utilis/component/leftsidebar";
import { Navbar } from "@/utilis/component/navbar";
import Link from "next/link";

export default function Groups() {
    const [isMobileRightSidebarOpen, setIsMobileRightSidebarOpen] = useState(false);
    const [groups, setGroups] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function getGroupData() {
            try {
                setIsLoading(true);
                const groupData = await fetchUserInfo("api/groups/getusergroups");
                setGroups(groupData || []);
            } catch (err) {
                setError("Failed to fetch groups");
                console.error("Error fetching groups:", err);
            } finally {
                setIsLoading(false);
            }
        }
        getGroupData();
    }, []);

    if (isLoading) {
        return (
            <div className="groups-hero">
                <Navbar setIsMobileRightSidebarOpen={setIsMobileRightSidebarOpen} />
                <Leftsidebar />
                {/* <Rightsidebar isMobileOpen={isMobileRightSidebarOpen} /> */}
                <ChatApplication/>
                <main className="groups-main">
                    <div className="flex items-center justify-center h-full">
                        <p>Loading groups...</p>
                    </div>
                </main>
            </div>
        );
    }

    if (error) {
        return (
            <div className="groups-hero">
                <Navbar setIsMobileRightSidebarOpen={setIsMobileRightSidebarOpen} />
                <Leftsidebar />
                {/* <Rightsidebar isMobileOpen={isMobileRightSidebarOpen} /> */}
                <ChatApplication/>
                <main className="groups-main">
                    <div className="flex items-center justify-center h-full">
                        <p className="text-red-500">{error}</p>
                    </div>
                </main>
            </div>
        );
    }

    return (
        <div className="groups-hero">
            <Navbar setIsMobileRightSidebarOpen={setIsMobileRightSidebarOpen} />
            <Leftsidebar />
            {/* <Rightsidebar isMobileOpen={isMobileRightSidebarOpen} /> */}
            <ChatApplication/>

            <main className="groups-main">
                <div className="groups-header">
                    <h1 className="groups-title">Your Groups</h1>
                    <p className="groups-subtitle">Join communities that match your interests</p>
                    <div className="create-group">
                        <CreateGroupForm />
                    </div>
                </div>

                <div className="groups-grid">
                    {groups.length === 0  || !groups ? (
                        <p>No groups found. Join or create a group to get started!</p>
                    ) : (
                        groups.map((group) => (
                            <div key={group.Id} className="group-card">
                                <div className="group-card-content">
                                    <div className="group-avatar-container">
                                        <img 
                                            src={
                                                group?.image
                                                  ? `/api/images?path=${group.image}`
                                                  : 'https://i.pinimg.com/736x/c1/e9/90/c1e990f02b655afa6bda4901bc1555f0.jpg'
                                              }
                                            alt={`${group.title} avatar`}
                                            className="group-avatar"
                                            onError={(e) => {
                                                e.target.onerror = null;
                                                e.target.src = 'https://i.pinimg.com/736x/c1/e9/90/c1e990f02b655afa6bda4901bc1555f0.jpg';
                                            }}
                                        />
                                    </div>
                                    
                                    <div className="group-info">
                                        <h3 className="group-name">{group.title}</h3>
                                        <p className="group-member-count">{group.nbr} members</p>
                                        <p className="group-description">{group.Descreption}</p>
                                    </div>
                                    
                                    <div className="group-action">
                                        <Link href={`/group/${group.Id}`} className="group-button visit-button">
                                            Visit
                                        </Link>
                                    </div>
                                </div>
                                <div className="group-card-glow"></div>
                            </div>
                        ))
                    )}
                </div>
            </main>
        </div>
    );
}
export function CreateGroupForm() {
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [selectedImage, setSelectedImage] = useState(null);
    const [formData, setFormData] = useState({
        title: '',
        description: '',
        image: null,
    });

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setSelectedImage(URL.createObjectURL(file));
            setFormData({ ...formData, image: file });
        }
    };

    const removeImage = () => {
        setSelectedImage(null);
        setFormData({ ...formData, image: null });
        const fileInput = document.getElementById('groupImage');
        if (fileInput) fileInput.value = '';
    };

    const handleInputChange = (e) => {
        const { id, value } = e.target;
        setFormData({ ...formData, [id]: value });
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        // Here you can handle the form submission
        console.log('Form Data:', formData);

        // Example: Send formData to an API
        const submitData = new FormData();
        submitData.append('title', formData.title);
        submitData.append('description', formData.description);
        if (formData.image) {
            submitData.append('image', formData.image);
        }

        // Example API call (replace with your actual API endpoint)
        fetch('/api/groups/add', {
            method: 'POST',
            body: submitData,
            credentials: 'include',
        })
        .then(response => response.json())
        .then(data => {
            console.log('Success:', data);
            // Handle success (e.g., show a success message, reset the form, etc.)
            resetForm();
        })
        .catch((error) => {
            console.error('Error fetching:', error);
            // Handle error (e.g., show an error message)
        });
    };

    const resetForm = () => {
        setFormData({
            title: '',
            description: '',
            image: null,
        });
        setSelectedImage(null);
        setShowCreateForm(false);
    };

    return (
        <div className="create-group">
            <div className="create-group-input">
                <div className="user-avatar">
                    <FontAwesomeIcon icon={faUserGroup} />
                </div>
                <button 
                    className="create-button"
                    onClick={() => setShowCreateForm(true)}
                >
                    <FontAwesomeIcon icon={faPlus}/>
                    Create Your Group
                </button>
            </div>

            {showCreateForm && (
                <div className="popup-overlay" onClick={() => setShowCreateForm(false)}>
                    <div className="popup-content" onClick={e => e.stopPropagation()}>
                        <button 
                            className="close-button"
                            onClick={() => setShowCreateForm(false)}
                        >
                            <FontAwesomeIcon icon={faTimes} />
                        </button>
                        
                        <h2>Create New Group</h2>
                        <form className="create-group-form" onSubmit={handleSubmit}>
                            <div className="form-group">
                                <input 
                                    type="text" 
                                    id="title" 
                                    placeholder="Enter group title"
                                    value={formData.title}
                                    onChange={handleInputChange}
                                />
                            </div>

                            <div className="form-group">
                                <textarea 
                                    id="description" 
                                    placeholder="Describe your group..."
                                    rows="4"
                                    value={formData.description}
                                    onChange={handleInputChange}
                                ></textarea>
                            </div>

                            <div className="form-group">
                                <div className="image-upload-container">
                                    <input 
                                        type="file" 
                                        id="groupImage" 
                                        accept="image/*"
                                        onChange={handleImageChange}
                                    />
                                    <label htmlFor="groupImage" className="image-upload-label">
                                        <FontAwesomeIcon icon={faImage} />
                                        <span>Choose Group Image</span>
                                    </label>
                                    {selectedImage && (
                                        <div className="image-preview">
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

                            <button type="submit" className="submit-button">
                                Create Group
                                <FontAwesomeIcon icon={faArrowRight} />
                            </button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}