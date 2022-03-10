import PageController from "./PageController";
import { BrowserRouter, Routes, Route} from "react-router-dom";
import About from "./Pages/About"
import Layout from "./Pages/Layout";

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Layout/>}>
                    <Route index element={<PageController/>} />
                </Route>
                <Route path="/about" element={<About/>}/>
            </Routes>
        </BrowserRouter>
        // <div>
        //     <Profile/>
        //     <PageController/>
        //     <ForkMe/>
        // </div>
    )
}

export default App
