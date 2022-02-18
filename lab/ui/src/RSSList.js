import React, { Component } from 'react'
import $ from 'jquery'
import './RSSList.css'
import Episodes from "./Episodes"

class RSSList extends Component {
    constructor(props) {
        super(props);
        this.state = {
            rsses:[],
            rssName: "",
        };
    }

    componentDidMount() {
        fetch('/rsses')
            .then(res => res.json())
            .then(res => {
                this.setState({
                    rsses: res
                })
            })
    }

    loadEpisodes(name) {
        this.setState({rssName: name})
    }

    render() {
        const rsses = this.state.rsses
        const name = this.state.rssName

        let episodes
        if (name) {
            episodes = (
                <Episodes rss={name} key={name}/>
            )
        }
        return (
            <div>
                <div id="rss_list">
                    {rsses.map((item, index) => (
                        <div rss={item.metadata.name} key={index}>
                            <img width="150px" src={item.spec.image} onClick={() => this.loadEpisodes(item.metadata.name)}></img>
                        </div>
                    ))}
                </div>
                {episodes}
            </div>
        )
    }
}

export default RSSList
