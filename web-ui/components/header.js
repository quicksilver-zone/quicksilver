import React, { useState, useEffect } from "react";
import { Navbar, NavbarBrand } from "reactstrap";
import Image from "next/image";
import CNWeb from "@/public/logo/cnweb-30.png";
import Link from "next/link";

const navButton = [
    {
        text: "How It Works",
        path: "#how-it-works",
    },
    {
        text: "Project",
        path: "#project",
    },
    {
        text: "Github",
        path: "#github",
    },
];

const Header = ({ currentPath }) => {
    return (
        <Navbar
            className="header-container fixed w-full h-30 z-10"
            light
            expand="md"
        >
            <div className="logo flex items-center font-semibold text-xl">
                <Image
                    src={CNWeb}
                    alt="cnweb logo"
                    className="invert dark:invert-0 h-10 w-10"
                />
                <NavbarBrand href="/" className="text-black dark:text-white">
                    CNWeb-30
                </NavbarBrand>
            </div>
            <div className="nav-bar grid text-2xl text-center relative">
                {navButton.map((button, index) => {
                    return (
                        <Link
                            key={index}
                            href={button.path}
                            className="text-xl leading-10 text-black dark:text-white font-thin hover:text-red dark:hover:text-red"
                            style={{
                                color: currentPath === button.text && "#C4181A",
                            }}
                        >
                            {button.text}
                        </Link>
                    );
                })}
            </div>
            <div className="header-right flex items-center">
                <Link className="get-started leading-10 font-thin text-white w-32 items-center text-center" href="/login">
                    Get Started
                </Link>
            </div>

        </Navbar>
    );
};

export default Header;
