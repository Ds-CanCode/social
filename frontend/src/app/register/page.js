"use client";
import { useState, useEffect } from 'react';
import './register.css';
import { LoginButton } from '../login/page.js';
import { Footer, Navbarrend } from '../page.js';
import { useRouter } from 'next/navigation';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { config } from "@fortawesome/fontawesome-svg-core";
import "@fortawesome/fontawesome-svg-core/styles.css";

config.autoAddCss = false; // Disable automatic CSS injection

import { 
  faUser, 
  faEnvelope, 
  faLock, 
  faCalendar, 
  faInfoCircle, 
  faImage,
  faExclamationTriangle,
  faSpinner
} from '@fortawesome/free-solid-svg-icons';

export default function Register() {
  const router = useRouter();

  const [formData, setFormData] = useState({
    email: '',
    password: '',
    firstName: '',
    lastName: '',
    aboutMe: '',
    dateOfBirth: '',
    avatar: null,
    nickName: ''
  });

  const [progress, setProgress] = useState(0);
  const [isRequired, setIsRequired] = useState(false);
  const [error, setError] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const requiredFields = ['email', 'password', 'firstName', 'lastName'];
    const filledFields = requiredFields.filter(field => formData[field].trim() !== '');
    const newProgress = (filledFields.length / requiredFields.length) * 100;
    setProgress(newProgress);
    setIsRequired(newProgress === 100);
  }, [formData]);

  const handleChange = (e) => {
    const { name, value, files } = e.target;
    if (name === 'avatar' && files) {
      setFormData(prev => ({
        ...prev,
        avatar: files[0]
      }));
    } else {
      setFormData(prev => ({
        ...prev,
        [name]: value
      }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);

    if (formData.email === '' || formData.password === '' || formData.firstName === '' || formData.lastName === '' || formData.dateOfBirth === '' ) { 
      setError("Please fill out all required fields");
      setIsSubmitting(false);
      return;
    } 

    try {
      const formDataToSend = new FormData();

      formDataToSend.append("email", formData.email);
      formDataToSend.append("password", formData.password);
      formDataToSend.append("firstName", formData.firstName);
      formDataToSend.append("lastName", formData.lastName);
      formDataToSend.append("dateOfBirth", formData.dateOfBirth);
      if (formData.nickName) {
        formDataToSend.append("nickName", formData.nickName);
      } 
      if (formData.aboutMe) {
        formDataToSend.append("aboutMe", formData.aboutMe);
      }
      if (formData.avatar) {
        formDataToSend.append("avatar", formData.avatar);
      }
            
      const response = await fetch('/api/register', {
        method: 'POST',
        body: formDataToSend
      });

      if (response.status === 201) {
        router.push("/login");
      } else {
        const data = await response.json();
        setError(data.error);
      }
    } catch (error) {
      console.log(error.message);
      setError("Registration failed please try again later");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className='register-hero'>
      <Navbarrend NavButton={() => <LoginButton text="Login" path="/login" />} />
      <div className="register-container">
        <div className="background-shapes">
          <div className="shape shape-1"></div>
          <div className="shape shape-2"></div>
        </div>

        <div className="background-glow"></div>
        
        <div className="register-card">
          <div className="progress-bar">
            <div className="progress-fill" style={{ width: `${progress}%` }}></div>
          </div>

          <div className="register-header">
            <h2>Create Account</h2>
            <p>Join our community</p>
          </div>

          <form onSubmit={handleSubmit} className="register-form">
            {error && (
              <div className='error-msg'>
                <FontAwesomeIcon icon={faExclamationTriangle} className="error-icon" />
                {error}
              </div>
            )}
            
            <div className="form-section required">
              <p className='info-lable'>
                <FontAwesomeIcon icon={faInfoCircle} className="info-icon" />
                required info
              </p>
              <div className="input-grid">
                <div className="input-wrapper">
                  <div className="input-icon">
                    <FontAwesomeIcon icon={faUser} />
                  </div>
                  <input
                    type="text"
                    name="firstName"
                    placeholder="First Name"
                    id="firstname"
                    required
                    value={formData.firstName}
                    onChange={handleChange}
                  />
                  <div className="input-glow"></div>
                </div>

                <div className="input-wrapper">
                  <div className="input-icon">
                    <FontAwesomeIcon icon={faUser} />
                  </div>
                  <input
                    type="text"
                    name="lastName"
                    id="lastname"
                    placeholder="Last Name"
                    required
                    value={formData.lastName}
                    onChange={handleChange}
                  />
                  <div className="input-glow"></div>
                </div>
              </div>

              <div className="input-wrapper">
                <div className="input-icon">
                  <FontAwesomeIcon icon={faEnvelope} />
                </div>
                <input
                  type="email"
                  name="email"
                  id="email"
                  placeholder="Email"
                  required
                  value={formData.email}
                  onChange={handleChange}
                />
                <div className="input-glow"></div>
              </div>

              <div className="input-wrapper">
                <div className="input-icon">
                  <FontAwesomeIcon icon={faLock} />
                </div>
                <input
                  type="password"
                  name="password"
                  id='password'
                  placeholder="Password"
                  required
                  value={formData.password}
                  onChange={handleChange}
                />
                <div className="input-glow"></div>
              </div>

              <div className="input-wrapper">
                <div className="input-icon">
                  <FontAwesomeIcon icon={faCalendar} />
                </div>
                <input
                  type="date"
                  name="dateOfBirth"
                  value={formData.dateOfBirth}
                  onChange={handleChange}
                  required
                />
                <div className="input-glow"></div>
              </div>
            </div>

            <div className={`form-section optional ${isRequired ? 'revealed' : ''}`}>
              <p className='info-lable'>
                <FontAwesomeIcon icon={faInfoCircle} className="info-icon" />
                optional info
              </p>
              <div className="input-wrapper">
                <div className="input-icon">
                  <FontAwesomeIcon icon={faUser} />
                </div>
                <input
                  type="text"
                  name="nickName"
                  id="nickName"
                  placeholder="Nick Name"
                  value={formData.nickName}
                  onChange={handleChange}
                />
                <div className="input-glow"></div>
              </div>

              <div className="input-wrapper">
                <textarea
                  name="aboutMe"
                  placeholder="About Me"
                  rows="3"
                  value={formData.aboutMe}
                  onChange={handleChange}
                ></textarea>
                <div className="input-glow"></div>
              </div>

              <div className="file-input-wrapper">
                <label>
                  <FontAwesomeIcon icon={faImage} className="file-icon" />
                  Profile Picture
                </label>
                <input
                  type="file"
                  name="avatar"
                  accept="image/*"
                  onChange={handleChange}
                />
              </div>
            </div>

            <button 
              type="submit" 
              id="register-button" 
              className="submit-button"
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <>
                  <FontAwesomeIcon icon={faSpinner} spin className="spinner-icon" />
                  Processing...
                </>
              ) : (
                'Join Now'
              )}
              <div className="button-glow"></div>
            </button>
          </form>
        </div>
      </div>
      <Footer/>
    </div>
  );
}