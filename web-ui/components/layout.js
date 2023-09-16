import React, { useEffect } from "react";
import Head from "next/head";
import style from '@/styles/Home.module.css'

const Layout = (props) => {
  return (
    <div>
      <Head>
        <title>{props.pageTitle}</title>
      </Head>
      <div class={`${style.bg}`}>
        <img src="/imgs/background.jpg" alt="" />
      </div>
      <div>
        {props.children}
      </div>
    </div>
  );
};
export default Layout;

