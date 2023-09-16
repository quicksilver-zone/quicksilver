import React, { useState, useEffect } from "react";
import { Navbar, NavbarBrand } from "reactstrap";
import { Box, Image as ChakraImage, Flex } from "@chakra-ui/react";
import Link from "next/link";
import { Center } from "@chakra-ui/react";
import { usePathname } from 'next/navigation'

const navButton = [
    {
        text: "Staking",
        path: "/staking",
        img: '/icons/staking.svg',
    },
    {
        text: "Assets",
        path: "/assets",
        img: '/icons/assets.svg',
    },
    {
        text: "Defi",
        path: "/defi",
        img: '/icons/defi.svg',
    },
    {
        text: "Airdrop",
        path: "/airdrop",
        img: '/icons/airdrop.svg',
    },
    {
        text: "Governance",
        path: "/governance",
        img: '/icons/governance.svg',
    }
];

const SideBar = () => {
    const currentPath = usePathname()
    const sidebar = React.createRef();
    const [width, setWidth] = useState(233);
    const [full, setFull] = useState(true);
    const [isShows, setIsShows] = useState([]);

    useEffect(() => {
        let showArray = navButton.map((_) => {
            return false;
        });
        setIsShows([...showArray]);
    }, []);

    useEffect(() => {
        const updateWidth = () => {
            const container = document.querySelector('.sidebar');
            if (container) {
                setWidth(container.offsetWidth);
            }
        };

        updateWidth();
        window.addEventListener('resize', updateWidth);

        return () => {
            window.removeEventListener('resize', updateWidth);
        };
    }, []);

    useEffect(() => {
        if (width < 233) {
            console.log(width)
            setFull(false);
        }
    }, [width]);

    const resize = () => {
        setFull(!full);
    }

    const handleMouseEnter = (button, index) => {
        if (button.subNav) {
            let showArray = [...isShows];
            showArray[index] = true;
            setIsShows([...showArray]);
        }
    };

    const handleMouseLeave = (button, index) => {
        if (button.subNav) {
            let showArray = [...isShows];
            showArray[index] = false;
            setIsShows([...showArray]);
        }
    };

    return (
        <Center margin={'20px'} zIndex={1}>
            <div className={full ? `sidebar` : `sidebar minimal-size`} ref={sidebar}>
                <Center>
                    <Navbar
                        className="menu"
                        light
                        expand="md"
                    >
                        <Flex justify={'space-between'} direction={'column'}>
                            <NavbarBrand className="text-black dark:text-white logo flex items-center font-semibold text-xl">
                                {
                                    full ? <ChakraImage alt="cnweb logo" src={'/logo/qs-text.svg'} />
                                        : <ChakraImage alt="cnweb logo" src={'/logo/qs_logo.svg'} boxSize={'100%'} />
                                }
                                <button 
                                    className={full ? `resize-btn` : `resize-btn minimal-btn`} onClick={resize} style={{left: full ? '2.3em' : '-.9em'}}>
                                    <span className="up-arrow"></span>
                                    <span className="down-arrow"></span>
                                </button>
                            </NavbarBrand>
                            <div className="menu-bar">
                                {navButton.map((button, index) => {
                                    return (
                                        <Link
                                            href={button.path}
                                            key={index}
                                            onMouseEnter={() => {
                                                handleMouseEnter(button, index);
                                            }}
                                            onMouseLeave={() => {
                                                handleMouseLeave(button, index);
                                            }}
                                            style={{
                                                background: currentPath.includes(button.path) && "rgba(231, 119, 40, 1)",
                                                boxShadow: currentPath.includes(button.path) && "7px 7px 22px #242424, -7px -7px 22px #383838",
                                            }}

                                        >
                                            <ChakraImage src={button.img} boxSize={'80%'} />
                                            <p style={{
                                                color: 'rgba(255, 255, 255, 1)',
                                                display: !full && 'none',
                                                color: currentPath.includes(button.path)? 'rgba(14, 14, 14, 1)' : 'rgba(255, 255, 255, 1)'
                                            }}
                                            >
                                                {button.text}
                                            </p>
                                        </Link>
                                    );
                                })}
                            </div>
                            {full ? <div className="menu-down flex items-center" style={{ color: '#979797', height:'18px', fontSize: '12px'}}>
                                Powered by Quicksilver Protocol.
                            </div> : <Box h='18px'/> } 
                        </Flex>
                    </Navbar>
                </Center>
            </div>
        </Center>
    );
};

export default SideBar;