import React, {Component} from "react";
import Detail from "./Detail";
import RSSList from "./RSSList";
import EpisodeItem from "./EpisodeItem";

class PageController extends Component {
    constructor(props) {
        super(props);
        this.state = {
            page: 'list', // supported value: list, detail, episode
            rss: '',
            episode: '',
        }
    }

    goPage(page, rss, episode) {
        this.setState({
            page: page,
            rss: rss,
            episode: episode,
        })
    }

    render() {
        const page = this.state['page']
        const rss = this.state['rss']
        const episode = this.state['episode']
        if (page === 'detail') {
            return (
                <Detail goHome={() => this.goPage('list')}
                        goEpisode={(rss, name) => this.goPage('episode', rss, name)} name={rss}/>
            )
        } else if (page === 'list') {
            return (
                <RSSList goDetail={(page, rss) => this.goPage(page, rss)}/>
            )
        } else if (page === 'episode') {
            return (
                <EpisodeItem goHome={() => this.goPage('list')}
                             goRSS={(rss) => this.goPage('detail', rss)} name={episode} rss={rss} />
            )
        } else {
            return (
                <div>Not found</div>
            )
        }
    }
}

export default PageController
