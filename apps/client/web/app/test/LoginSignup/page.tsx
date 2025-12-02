"use client"

import React from 'react'
import './LoginSignup.css'
import Image from 'next/image'
import user_icon from '../../../Assets/user.png'
import email_icon from '../../../Assets/mail.png'
import password_icon from '../../../Assets/Password.png'
import { useState } from 'react'

const LoginSignup: React.FC = () => {

    const [action,setAction] = useState<"Login" | "Sign Up">("Login");

    return(
        <div className = 'container'>
            <div className = "header">
                <div className = "text">{action}</div>
                <div className = "underline"></div>
            </div>
            <div className = "inputs"> 
                {action==="Login"?<div></div>:<div className = "input">
                    <div className="input-icon">
                        <Image src={user_icon} alt="User icon" width={24} height={24} />
                    </div>
                    <input type = "text" placeholder = "Name"/>
                </div>}

                <div className = "input">
                    <div className="input-icon">
                        <Image src={email_icon} alt="Email icon" width={24} height={24} />
                    </div>
                    <input type = "email" placeholder= "Email Id"/>
                </div>
                <div className = "input">
                    <div className="input-icon">
                        <Image src={password_icon} alt="Password icon" width={24} height={24} />
                    </div>
                    <input type = "password" placeholder = "Password"/>
                </div>
            </div>
            {action==="Sign Up"?<div></div>: <div className="forgot-password">Lost Password? <span>Click Here!</span></div>}
            <div className = "submit-container">
                <div className = {action==="Login"?"submit gray":"submit"} onClick = {()=>{setAction("Sign Up")}}>Sign Up</div>
                <div className = {action==="Sign Up"?"submit gray":"submit"}onClick = {()=>{setAction("Login")}}>Login</div>
            </div>
        </div>
    )
}

export default LoginSignup 