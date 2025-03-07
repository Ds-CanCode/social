"use client";
// import React, { useState } from 'react';
import "./login.css";
import { Navbarrend, Footer } from "../page.js";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useRouter } from 'next/navigation';
import { checkAuth } from "@/utilis/ auth"; 
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";

import { 
  faEnvelope, 
  faLock, 
} from '@fortawesome/free-solid-svg-icons';

export function LoginButton({ text, path }) {
  return (
    <Link href={path} className="login-buttonn">
      <span className="button-glow"></span>
      {text}
    </Link>
  );
}

export default function Login() {
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true); // Add a loading state


  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');


    // Check if the user is already authenticated
    useEffect( () => {
      const verifyAuth = async () => {
        const isAuthenticated = await checkAuth();
        if (isAuthenticated) {
          router.push("/"); // Redirect to the index page if authenticated
        } else {
          setLoading(false); // Update the loading state
        }
      };
     verifyAuth();
    }, [router]);



  const handleSubmit = async (e) => {
    e.preventDefault();
    const userData = {
        email,
        password,
      };
      
      console.log(userData);
    try {
      const response = await fetch("/api/login", {
        method: "POST",
        headers: {
        },
        body: JSON.stringify(userData), credentials: "include"
      });

      
      if (response.status === 200) {
        router.push("/");
      } else {
        const data = await response.json();
        setError(data.error);
      }
    } catch (error) {
      console.error("Error during login:", error);
    }
  };

  if (loading) {
    return <div></div>
  }

  return (
    <div className="hero">
      <Navbarrend
        NavButton={() => <LoginButton text="Register" path="/register" />}
      />

      <div className="login-container">
        <div className="background-shapes">
          <div className="shape shape-1"></div>
          <div className="shape shape-2"></div>
        </div>

        <div className="login-card">
          <div className="login-header">
            <div className="avatar-circle">
              <svg
                width="45"
                height="45"
                viewBox="0 0 45 45"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M43.6189 22.2351L38.9828 22.1994M38.7145 13.8039L35.5461 18.5019C35.5278 18.529 35.5156 18.5597 35.5104 18.5919L35.2037 20.4836C35.1935 20.5462 35.2102 20.6102 35.2496 20.6599L36.3881 22.0954C36.4305 22.1489 36.4948 22.1803 36.5631 22.1808L38.9828 22.1994M38.7145 13.8039C38.5042 13.2782 37.9028 12.2511 37.1207 12.1576C37.0492 12.149 36.9811 12.1852 36.9386 12.2434L36.9186 12.2709M38.7145 13.8039L41.3139 13.8539M33.3547 17.1596L36.9186 12.2709M36.9186 12.2709C36.3638 11.8398 34.9943 11.1878 33.9213 11.9951C33.9028 12.009 33.8866 12.0263 33.8735 12.0454L33.8164 12.1285M31.3082 15.7863L33.8164 12.1285M33.8164 12.1285C33.4149 11.7892 32.3723 11.2827 31.3893 11.9441C31.3714 11.9561 31.3552 11.9713 31.3417 11.9881L28.9656 14.9516L30.8076 12.4011C30.9154 12.2519 30.8088 12.0435 30.6248 12.0435H29.1478C29.0994 12.0435 29.0523 12.059 29.0135 12.0878L26.8213 13.7126C26.797 13.7306 26.7766 13.7533 26.7612 13.7794L24.6311 17.4006C24.6155 17.4271 24.6055 17.4567 24.6017 17.4873L24.0077 22.2929C24.0035 22.3273 23.9913 22.3604 23.9722 22.3893L22.515 24.6003C22.4514 24.6967 22.471 24.8255 22.5602 24.8988L26.0663 27.777C26.0855 27.7927 26.1019 27.8115 26.115 27.8326L28.0604 30.9751C28.1466 31.1145 28.348 31.1182 28.4394 30.9822L30.5821 27.793C30.603 27.7619 30.6313 27.7365 30.6645 27.7191L35.0023 25.4418C35.0157 25.4347 35.0283 25.4264 35.04 25.4168L38.9828 22.1994"
                  stroke="white"
                  strokeWidth="1.5"
                />
                <path
                  d="M11.1757 5.27882L13.4628 9.3116M6.32627 13.7418L11.9791 14.1367C12.0116 14.139 12.0443 14.1341 12.0748 14.1225L13.8665 13.4424C13.9258 13.4198 13.9729 13.3734 13.9962 13.3145L14.6701 11.6107C14.6952 11.5473 14.6903 11.4758 14.6566 11.4165L13.4628 9.3116M6.32627 13.7418C5.97616 14.1867 5.38733 15.2211 5.6974 15.9452C5.72577 16.0114 5.7911 16.0523 5.86277 16.0599L5.8966 16.0635M6.32627 13.7418L4.73154 11.0145M11.9123 16.7056L5.8966 16.0635M5.8966 16.0635C5.8007 16.7596 5.9208 18.2716 7.15641 18.7972C7.17771 18.8063 7.20077 18.8117 7.22385 18.8135L7.3244 18.8213M11.7463 19.1645L7.3244 18.8213M7.3244 18.8213C7.23126 19.3387 7.31397 20.4948 8.37822 21.0154C8.39764 21.025 8.41881 21.0314 8.44018 21.0347L12.1947 21.6107L9.06483 21.2907C8.8818 21.272 8.75457 21.4685 8.84657 21.6278L9.58507 22.907C9.60926 22.9489 9.64626 22.9819 9.69063 23.0011L12.1938 24.0872C12.2215 24.0992 12.2514 24.1056 12.2817 24.1058L16.4828 24.14C16.5136 24.1402 16.5442 24.1341 16.5726 24.1221L21.0313 22.2337C21.0633 22.2201 21.098 22.2142 21.1326 22.2162L23.776 22.3727C23.8913 22.3795 23.9931 22.2982 24.0119 22.1843L24.7514 17.7088C24.7554 17.6844 24.7635 17.6607 24.7752 17.6389L26.5241 14.3828C26.6016 14.2385 26.5042 14.0622 26.3407 14.0511L22.5074 13.79C22.47 13.7875 22.4339 13.7757 22.4022 13.7556L18.2611 11.1377C18.2483 11.1296 18.2347 11.1228 18.2206 11.1175L13.4628 9.3116"
                  stroke="white"
                  strokeWidth="1.5"
                />
                <path
                  d="M12.5086 41.7006L14.8576 37.7035M22.2625 41.6688L19.7781 36.5759C19.7638 36.5466 19.7433 36.5207 19.718 36.5L18.2331 35.2885C18.1839 35.2484 18.1202 35.2309 18.0574 35.2401L16.245 35.5084C16.1775 35.5184 16.1181 35.5584 16.0836 35.6172L14.8576 37.7035M22.2625 41.6688C22.8229 41.7496 24.0131 41.7423 24.4851 41.1118C24.5283 41.0541 24.531 40.9771 24.5018 40.9112L24.488 40.8801M22.2625 41.6688L20.9195 43.895M22.0362 35.3493L24.488 40.8801M24.488 40.8801C25.1388 40.6151 26.3881 39.7551 26.2255 38.4222C26.2227 38.3992 26.2159 38.3766 26.2059 38.3557L26.1624 38.2647M24.2487 34.2636L26.1624 38.2647M26.1624 38.2647C26.657 38.0866 27.6169 37.4369 27.5357 36.255C27.5342 36.2334 27.5292 36.2118 27.5213 36.1917L26.143 32.6522L27.4308 35.5227C27.5061 35.6906 27.7399 35.7025 27.8319 35.5432L28.5704 34.264C28.5946 34.2221 28.6046 34.1736 28.5991 34.1255L28.2881 31.4147C28.2847 31.3847 28.2752 31.3556 28.2603 31.3293L26.1893 27.674C26.1741 27.6471 26.1536 27.6237 26.129 27.6051L22.2641 24.6879C22.2364 24.667 22.2139 24.64 22.1984 24.6089L21.0122 22.2414C20.9605 22.1382 20.8392 22.0907 20.7311 22.1314L16.4855 23.7287C16.4623 23.7374 16.4378 23.7422 16.413 23.743L12.7188 23.8565C12.555 23.8615 12.4511 24.034 12.5232 24.1812L14.2137 27.6314C14.2302 27.6651 14.238 27.7023 14.2365 27.7397L14.0399 32.635C14.0393 32.6501 14.0402 32.6653 14.0426 32.6802L14.8576 37.7035"
                  stroke="white"
                  strokeWidth="1.5"
                />
                <circle
                  cx="22.2994"
                  cy="22.9844"
                  r="21.3118"
                  stroke="white"
                  strokeWidth="1.5"
                />
              </svg>
            </div>
            <h1>Welcome Back</h1>
            <p>Please sign in to continue</p>
          </div>
          {error && <div className='error-msg'>{error}</div>}


          <form onSubmit={handleSubmit} className="login-form">
            <div className="input-group">
              <div className="input-wrapper">
              <FontAwesomeIcon className="input-icon" icon={faEnvelope}></FontAwesomeIcon>
                <input
                  type="email"
                  placeholder="Email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                />
              </div>
            </div>

            <div className="input-group">
              <div className="input-wrapper">
              <FontAwesomeIcon className="input-icon" icon={faLock}></FontAwesomeIcon>
                <input
                  type="password"
                  placeholder="Password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
              </div>
            </div>

            <button type="submit" className="login-button">
              Sign In
            </button>

            <Link href="/register" className="forgot-password">
              sign up
            </Link>
          </form>
        </div>
      </div>
      <Footer/>
    </div>
  );
}
