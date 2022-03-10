import ForkMe from "./ForkMe";
import { Outlet, Link } from "react-router-dom";
import Profile from "./Profile";
import BackTop from 'cuke-ui/lib/back-top';

export default function () {
    return (
        <>
            <ForkMe/>
            <Profile/>
            <BackTop visibilityHeight={300} style={{right: 50,bottom: 100}} />
            <Outlet/>
        </>
    )
}