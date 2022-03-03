import ForkMe from "./ForkMe";
import Profile from "./Profile";
import PageController from "./PageController";
import BackTop from 'cuke-ui/lib/back-top';

function App() {
    return (
        <div>
            <Profile/>
            <PageController/>
            <ForkMe/>

            <BackTop visibilityHeight={300} style={{right: 50,bottom: 100}} />
        </div>
    )
}

export default App
