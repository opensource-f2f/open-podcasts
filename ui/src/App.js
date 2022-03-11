import { BrowserRouter, Routes, Route} from "react-router-dom";
import About from "./Pages/About"
import Layout from "./Pages/Layout";
import RSSByCategory from "./Pages/RSSByCategory";
import Detail from "./Pages/Detail";
import EpisodeItem from "./Pages/EpisodeItem";

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Layout/>}>
                    <Route index element={<RSSByCategory/>} />
                    <Route path="/rsses/search" element={<RSSByCategory/>} />
                    <Route path="/rsses/:id/episodes" element={<Detail/>} />
                    <Route path="/rsses/:rss/episodes/:episode" element={<EpisodeItem/>}/>
                </Route>
                <Route path="/about" element={<About/>}/>
            </Routes>
        </BrowserRouter>
    )
}

export default App
