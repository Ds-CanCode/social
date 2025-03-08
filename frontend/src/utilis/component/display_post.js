"use client";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faThumbsUp,
  faThumbsDown,
  faComment,
  faClock,
  faPaperPlane,
  faImage,
} from "@fortawesome/free-solid-svg-icons";
import { fetchUserInfo } from "@/utilis/fetching_data.js";
import { useEffect, useState, useRef, use } from "react";
import "./css/display_post.css";
import Link from "next/link";

// Throttle function to limit function calls
function throttle(func, delay) {
  let lastCall = 0;
  return function(...args) {
    const now = new Date().getTime();
    if (now - lastCall < delay) {
      return;
    }
    lastCall = now;
    return func(...args);
  };
}

//****************** Create Post Component: a function that displays a single post ***********************/

export function Post({ post }) {
  const [isCommentSectionOpen, setIsCommentSectionOpen] = useState(false);
  const [comments, setComments] = useState([]);
  const [likes, setLikes] = useState(post.likes);
  const [dislikes, setDislikes] = useState(post.dislikes);
  const [isLoading, setIsLoading] = useState(false);

  

  const [image, setImage] = useState(null);
  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setImage(file);
      // setPreview(URL.createObjectURL(file));
    }
  };


  // const [comment, setNewComment] = useState("");

  const [formComments, setCommentData] = useState({
    postId: "",
    content: "",
  });

  const [userdata, setUserdata] = useState(null);
  useEffect(() => {
    async function getUserData() {
      const userdata = await fetchUserInfo("api/users/info");
      setUserdata(userdata); // Store the user data in state
    }
    getUserData();
  }, []);

  const handleSubmitComment = async (e) => {
    e.preventDefault();
    if (!formComments.content.trim() || isLoading) return;
    
    setIsLoading(true);
    const newCommentForm = new FormData();
    newCommentForm.append("post_id", formComments.postId);
    newCommentForm.append("content", formComments.content);
    if (image) newCommentForm.append("image", image);
    try {
      const response = await fetch("/api/comment/add", {
        method: "POST",
        credentials: "include",
        body: newCommentForm,
      });
      if (response.status === 201) {
        setCommentData({
          postId: "",
          content: "",
        });
        setImage(null);
      } else {
        console.log(response.status);
      }

      if (!response.ok) {
        throw new Error("Failed to add comment");
      }

      const data = await response.json();
      console.log(data);

      setComments([data, ...comments]);
    } catch (error) {
      console.error("Error submitting comment:", error);
    } finally {
      setIsLoading(false);
    }
  };

  // Throttled reaction handlers
  const handleReact = throttle(async (Id, type) => {
    if (isLoading) return;
    
    setIsLoading(true);
    try {
      const response = await fetch(`/api/reactPost/add`, {
        method: "post",
        credentials: "include",
        body: JSON.stringify({ post_id: Id, reaction_type: type }),
      });

      if (response.status === 200) {
        const data = await response.json();
        setLikes(data.like_count);
        setDislikes(data.dislike_count);
      } else {
        console.log("Failed to like post");
      }
    } catch (error) {
      console.error("Error liking post:", error);
    } finally {
      setIsLoading(false);
    }
  }, 500); // 500ms throttle delay

  const handleCommentReact = throttle(async (Id, CommentId, type, current) => {
    if (isLoading) return;
    
    setIsLoading(true);
    try {
      const response = await fetch(`/api/reactComment/add`, {
        method: "post",
        credentials: "include",
        body: JSON.stringify({ post_id: Id, comment_id: CommentId, reaction_type: type }),
      });

      if (response.status === 200) {
        const updatedComments = comments.map(comment => {
          if (comment.id === CommentId) {
            let updatedLikes = comment.likes;
            let updatedDislikes = comment.dislikes;
            console.log(current);
            
            if (type === 'LIKE') {
              if (current == '') {
                updatedLikes += 1
              } else if (updatedLikes > 0 && current == 'LIKE') {
                updatedLikes -= 1
              } else if (updatedDislikes > 0 && current == 'DISLIKE') {
                updatedLikes += 1
                updatedDislikes -= 1
              }
            } else if (type === 'DISLIKE') {
              if (current == '') {
                updatedDislikes += 1
              } else if (updatedDislikes > 0 && current == 'DISLIKE') {
                updatedDislikes -= 1
              } else if (updatedLikes > 0 && current == 'LIKE'){
                updatedLikes -= 1
                updatedDislikes += 1
              }
            }
         
            return {
              ...comment,
              likes: updatedLikes,
              dislikes: updatedDislikes,
            };
          }
          return comment;
        });
        setComments(updatedComments);
      } else {
        console.log("Failed to like post");
      }
    } catch (error) {
      console.error("Error liking post:", error);
    } finally {
      setIsLoading(false);
    }
  }, 500); // 500ms throttle delay

  const getCurrent = throttle(async (commentId, type) => {
    if (isLoading) return;
    
    setIsLoading(true);
    try {
      const response = await fetch(`/api/reactComment/getCurrent?commentId=${commentId}`, {
        method: "post",
        credentials: "include",
      });

      if (response.status === 200) {
        const data = await response.json();
        if (type === "LIKE") {
          handleCommentReact(post.id, commentId, "LIKE", data.currentReact);
        } else if (type === "DISLIKE") {
          handleCommentReact(post.id, commentId, "DISLIKE", data.currentReact);
        }
      } else {
        console.log("Failed to like post");
      }
    } catch (error) {
      console.error("Error liking post:", error);
    } finally {
      setIsLoading(false);
    }
  }, 500); // 500ms throttle delay

  const getComments = throttle(async (postId) => {
    if (isLoading) return;
    
    setIsLoading(true);
    try {
      const data = await fetchUserInfo(`api/comment/get?PostId=${postId}`);
      if (data) {
        setComments(data);
      }
    } catch (error) {
      console.error("Error fetching comments:", error);
    } finally {
      setIsLoading(false);
    }
  }, 1000); // 1000ms throttle delay for comments loading

  return (
    <article className="post">
      <div className="post-header">
        <div className="post-user-info">
          <Link href={`/profile/${post.creator?.id}`}>
            <div className="user-avatar">
              <img
                src={
                  post.creator?.avatar
                    ? `/api/images?path=${post.creator.avatar}`
                    : "/default-img.jpg"
                }
                alt="User avatar"
                width={40}
                height={40}
                className="avatar-image"
              />
            </div>
          </Link>
          <div className="user-details">
            <Link href={`/profile/${post.creator?.id}`}>
              <h3 className="user-name">
                {post.creator?.first_name + " " + post.creator?.last_name}
              </h3>
            </Link>

            <span className="post-timestamp">
              <FontAwesomeIcon icon={faClock} />
              <time>{post?.created_at}</time>
            </span>
          </div>
        </div>
      </div>

      <div className="post-content">
        <h2 className="post-title">{post?.title || "loading title..."}</h2>
        <p className="post-text">{post?.content || "loading content..."}</p>
        {post?.image && (
          <div className="post-image">
            <img
              src={`/api/images?path=${post.image}`}
              alt="Post image"
              className="post-image01"
            />
          </div>
        )}
      </div>

      <div className="post-actions">
        <button
          onClick={() => handleReact(post.id, "LIKE")}
          className={`action-button ${isLoading ? "disabled" : ""}`}
          disabled={isLoading}
        >
          <FontAwesomeIcon icon={faThumbsUp} />
          <span>{likes}</span>
        </button>
        <button
          onClick={() => handleReact(post.id, "DISLIKE")}
          className={`action-button ${isLoading ? "disabled" : ""}`}
          disabled={isLoading}
        >
          <FontAwesomeIcon icon={faThumbsDown} />
          <span>{dislikes}</span>
        </button>
        <button
          className={`action-button ${isLoading ? "disabled" : ""}`}
          disabled={isLoading}
          onClick={() => {
            getComments(post.id);
            setIsCommentSectionOpen(!isCommentSectionOpen);
          }}
        >
          <FontAwesomeIcon icon={faComment} />
          <span>{comments.length != 0 ? comments.length : ""}</span>
        </button>
      </div>

      {isCommentSectionOpen && (
        <div className="comments-section">
          <form className="comment-form" onSubmit={handleSubmitComment}>
            <div className="user-avatar">
              <img
                src={`/api/images?path=${userdata?.avatar}`}
                alt="Your avatar"
                width={32}
                height={32}
                className="avatar-image"
              />
            </div>
            <input
              type="text"
              className="comment-input"
              placeholder="Write a comment..."
              value={formComments.content}
              onChange={(e) =>
                setCommentData({
                  postId: post.id,
                  content: e.target.value,
                })
              }
            />
            <input
              className="comment-input-file"
              type="file"
              id="postImage"
              onChange={handleImageChange}
              accept="image/*"
              disabled={isLoading}
            />
            <label htmlFor="postImage" className="comment-input-file-label">
              <FontAwesomeIcon icon={faImage} />
            </label>
            <button
              type="submit"
              className="comment-submit"
              disabled={!formComments.content.trim() || isLoading}
            >
              <FontAwesomeIcon icon={faPaperPlane} />
            </button>
          </form>

          <div className="comments-list">
            {comments.map((comment) => (
              <div key={comment.id} className="comment">
                <div className="comment-header">
                  <div className="user-avatar">
                    <img
                      src={`/api/images?path=${comment.user?.avatar}`}
                      alt="Commenter avatar"
                      width={32}
                      height={32}
                      className="avatar-image"
                    />
                  </div>
                  <div className="comment-details">
                    <span className="comment-author">
                      {comment.user?.first_name + " " + comment.user?.last_name}
                    </span>
                    <span className="comment-timestamp">
                      <FontAwesomeIcon icon={faClock} />
                      <time>{comment.created_at}</time>
                    </span>
                  </div>
                </div>
                <div className="comment-divider">
                  <p className="comment-content">{comment?.content}</p>
                  {comment?.pathimg && (
                    <div className="comment-image">
                      <img
                        src={`/api/images?path=${comment.pathimg}`}
                        alt="Comment image"
                        className="comment-image01"
                      />
                    </div>
                  )}
                </div>
                <div className="comment-actions">
                  <button 
                    className={`action-button ${isLoading ? "disabled" : ""}`}
                    disabled={isLoading}
                    onClick={() => {
                      getCurrent(comment.id, "LIKE");
                    }}
                  >
                    <FontAwesomeIcon icon={faThumbsUp} />
                    <span>{comment.likes}</span>
                  </button>
                  <button 
                    className={`action-button ${isLoading ? "disabled" : ""}`}
                    disabled={isLoading}
                    onClick={() => {
                      getCurrent(comment.id, "DISLIKE");
                    }}
                  >
                    <FontAwesomeIcon icon={faThumbsDown} />
                    <span>{comment.dislikes}</span>
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </article>
  );
}

// ************** Post list: a function that displays a list of posts ****************//
export function PostList() {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [isFetching, setIsFetching] = useState(false);
  const loadMoreRef = useRef(null);

  // Throttled fetch posts function
  const fetchPosts = throttle(async (pageNum) => {
    if (isFetching || !hasMore) return;
    
    setIsFetching(true);
    try {
      const data = await fetchUserInfo(`api/post/getAll?offset=${pageNum}`);
      if (data && data.length > 0) {
        if (pageNum === 1) {
          setPosts(data);
        } else {
          setPosts(prevPosts => [...prevPosts, ...data]);
        }
        setPage(pageNum + 1);
      } else {
        setHasMore(false);
      }
    } catch (error) {
      setError("Failed to fetch posts");
    } finally {
      setLoading(false);
      setIsFetching(false);
    }
  }, 1000); // 1000ms throttle delay

  // Initial load
  useEffect(() => {
    fetchPosts(1);
  }, []);

  // Infinite scroll setup
  useEffect(() => {
    const observer = new IntersectionObserver(
      entries => {
        if (entries[0].isIntersecting && hasMore && !isFetching) {
          fetchPosts(page);
        }
      },
      { threshold: 0.5 }
    );
    
    if (loadMoreRef.current) {
      observer.observe(loadMoreRef.current);
    }
    
    return () => {
      if (loadMoreRef.current) {
        observer.unobserve(loadMoreRef.current);
      }
    };
  }, [hasMore, isFetching, page]);

  if (loading && posts.length === 0) return <p>Loading posts...</p>;
  if (error) return <p>{error}</p>;

  return (
    <div className="post-list">
      {posts.length > 0 ? (
        <>
          {posts.map((post) => <Post key={post.id} post={post} />)}
          <div ref={loadMoreRef} className="load-more">
            {isFetching && <p>Loading more posts...</p>}
          </div>
        </>
      ) : (
        <p>No posts available.</p>
      )}
    </div>
  );
}

// Add CSS styles for disabled buttons
const styleElement = typeof document !== 'undefined' ? document.createElement('style') : null;
if (styleElement) {
  styleElement.textContent = `
    .action-button.disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
    
    .load-more {
      text-align: center;
      padding: 1rem;
      margin-top: 1rem;
    }
  `;
  document.head.appendChild(styleElement);
}

// ******** Post container: a function that displays a form to create a new post **********//
export function PostContainer() {
  const [imageFile, setImageFile] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [formData, setFormData] = useState({
    title: "",
    content: "",
    privacy: "public",
  });
  const [imagePreview, setImagePreview] = useState(null);
  const modalRef = useRef(null);
  const fileInputRef = useRef(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const [targetedFriends, setTargetedFriends] = useState(null);

  // Add these state variables at the top of your PostContainer component
  const [showFriendsModal, setShowFriendsModal] = useState(false);
  const [selectedFriends, setSelectedFriends] = useState([]);
  const friendsModalRef = useRef(null);

  // Add this function to handle friend selection
  const handleFriendSelection = (friendId) => {
    setSelectedFriends((prev) => {
      if (prev.includes(friendId)) {
        return prev.filter((id) => id !== friendId);
      }
      return [...prev, friendId];
    });
  };

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (modalRef.current && !modalRef.current.contains(event.target)) {
        setIsModalOpen(false);
      }
    };

    if (isModalOpen) {
      document.addEventListener("mousedown", handleClickOutside);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [isModalOpen]);

  const handleImageUpload = (e) => {
    const file = e.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result);
      };
      reader.readAsDataURL(file);
      setImageFile(file);
    }
  };

  const removeImage = () => {
    setImagePreview(null);
    setImageFile(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  // Throttled fetch friends function
  const fetchFriendshandler = throttle(async () => {
    try {
      const response = await fetchUserInfo(
        `api/users/userfollowers?profileId=0`
      );
      setTargetedFriends(response);
    } catch (error) {
      console.error("Error fetching friends:", error);
    }
  }, 1000);

  // Throttled post submission
  const handleSubmit = throttle(async (e) => {
    e.preventDefault();
    if (isSubmitting) return;

    setIsSubmitting(true);
    try {
      const newPostForm = new FormData();

      newPostForm.append("title", formData.title);
      newPostForm.append("content", formData.content);
      newPostForm.append("privacy", formData.privacy);
      if (imageFile) newPostForm.append("image", imageFile);
      if (selectedFriends.length > 0) {
        newPostForm.append("targetedFriends", selectedFriends);
        newPostForm.set("privacy", "semi-private");
      }

      const response = await fetch("/api/post/add", {
        method: "POST",
        body: newPostForm,
        credentials: "include",
      });

      if (response.status === 201) {
        console.log("Post created successfully");
        setIsModalOpen(false);
        setFormData({
          title: "",
          content: "",
          privacy: "public",
        });
        setImagePreview(null);
        setSelectedFriends([]);
      } else {
        console.log(response.status);
      }
    } catch (error) {
      console.error("Error during post creation:", error);
    } finally {
      setIsSubmitting(false);
    }
  }, 1000);

  useEffect(() => {
    fetchFriendshandler();
  }, []);

  return (
    <div className="content-container">
      <div className="create-post-trigger" onClick={() => setIsModalOpen(true)}>
        <div className="trigger-input">What's on your mind?</div>
      </div>

      {isModalOpen && ( 
        <div className="modal-overlay">
          <div className="create-post-modal" ref={modalRef}>
            <div className="modal-header">
              <h2>Create Post</h2>
              <button
                className="close-button"
                onClick={() => setIsModalOpen(false)}
                aria-label="Close modal"
              >
                ×
              </button>
            </div>

            <form onSubmit={handleSubmit} className="post-form">
              <div className="form-group">
                <div className="image-upload-containers">
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleImageUpload}
                    className="file-input"
                    ref={fileInputRef}
                    id="image-upload"
                    disabled={isSubmitting}
                  />
                  <label htmlFor="image-upload" className={`upload-labels ${isSubmitting ? "disabled" : ""}`}>
                    <svg
                      width="24"
                      height="24"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    >
                      <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h7" />
                      <path d="M16 5V3" />
                      <path d="M8 14a3 3 0 0 1 3-3h4" />
                      <line x1="17" y1="8" x2="17" y2="14" />
                      <line x1="14" y1="11" x2="20" y2="11" />
                    </svg>
                    Upload Image
                  </label>
                </div>
                {imagePreview && (
                  <div className="image-preview-container">
                    <img
                      src={imagePreview}
                      alt="Preview"
                      className="image-preview"
                    />
                    <button
                      type="button"
                      className="remove-image"
                      onClick={removeImage}
                      disabled={isSubmitting}
                    >
                      ×
                    </button>
                  </div>
                )}
              </div>

              <div className="form-group">
                <input
                  type="text"
                  placeholder="Post title"
                  value={formData.title}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      title: e.target.value,
                    })
                  }
                  className="form-input"
                  required
                  disabled={isSubmitting}
                />
              </div>

              <div className="form-group">
                <textarea
                  placeholder="Write your post..."
                  value={formData.content}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      content: e.target.value,
                    })
                  }
                  className="form-textarea"
                  required
                  disabled={isSubmitting}
                />
              </div>

              <div className="privacy-options">
                <label className="radio-label">
                  <input
                    type="radio"
                    name="privacy"
                    value="public"
                    checked={formData.privacy === "public"}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        privacy: e.target.value,
                      })
                    }
                    disabled={isSubmitting}
                  />
                  <span className="radio-custom"></span>
                  Public
                </label>
                <label className="radio-label">
                  <input
                    type="radio"
                    name="privacy"
                    value="private"
                    checked={formData.privacy === "private"}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        privacy: e.target.value,
                      })
                    }
                    disabled={isSubmitting}
                  />
                  <span className="radio-custom"></span>
                  Private
                </label>
              </div>
              <div className="target-friends-section">
                <button
                  type="button"
                  className={`target-friends-button ${isSubmitting ? "disabled" : ""}`}
                  onClick={() => {
                    fetchFriendshandler();
                    setShowFriendsModal(true);
                  }}
                  disabled={isSubmitting}
                >
                  Target Friends
                </button>
              </div>
              {showFriendsModal && (
                <div className="friends-modal-overlay">
                  <div className="friends-modal" ref={friendsModalRef}>
                    <div className="friends-modal-header">
                      <h3>Select Friends</h3>
                      <button
                        className="close-button"
                        onClick={() => setShowFriendsModal(false)}
                        aria-label="Close modal"
                      >
                        ×
                      </button>
                    </div>

                    <div className="friends-list">
                      {targetedFriends && targetedFriends.map((friend) => (
                        <label key={friend.Id} className="friend-item">
                          <input
                            type="checkbox"
                            checked={selectedFriends.includes(friend.Id)}
                            onChange={() => handleFriendSelection(friend.Id)}
                            disabled={isSubmitting}
                          />
                          <span className="friend-name">
                            {friend.FirstName + " " + friend.LastName}
                          </span>
                        </label>
                      ))}
                    </div>

                    <div className="friends-modal-footer">
                      <button
                        className={`confirm-selection ${isSubmitting ? "disabled" : ""}`}
                        onClick={() => {
                          console.log("Selected friends:", selectedFriends);
                          setShowFriendsModal(false);
                        }}
                        disabled={isSubmitting}
                      >
                        Confirm Selection
                      </button>
                    </div>
                  </div>
                </div>
              )}
              <button 
                type="submit" 
                className={`submit-button ${isSubmitting ? "disabled" : ""}`}
                disabled={isSubmitting}
              >
                {isSubmitting ? "Creating Post..." : "Create Post"}
                <div className="button-glow"></div>
              </button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}