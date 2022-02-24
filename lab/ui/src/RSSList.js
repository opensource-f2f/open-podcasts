import React, { Component } from 'react'
import $ from 'jquery'
import './RSSList.css'
import Episodes from "./Episodes"

class RSSItem extends Component {
    constructor(props) {
        super(props);
    }
    render() {
        if (this.props.image) {
            return (
                <div rss={this.props.name}>
                    <img width="150px" src={this.props.image} alt={this.props.name}
                         onClick={() => this.props.loadEpisodes(this.props.name)}/>
                </div>
            )
        }
        return ''
    }
}

class RSSList extends Component {
    constructor(props) {
        super(props);
        this.state = {
            rsses:[],
            rssName: "",
        };
        this.loadEpisodes = this.loadEpisodes.bind(this);
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
                        <RSSItem name={item.metadata.name} key={index} image={item.spec.image}
                                 loadEpisodes={this.loadEpisodes}/>
                    ))}
                </div>
                {episodes}
            </div>
        )
    }
}

export default RSSList
