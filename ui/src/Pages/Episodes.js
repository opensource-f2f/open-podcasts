import React, { Component } from 'react'
import './Episodes.css'
import $ from 'jquery'
import Button from 'cuke-ui/lib/button';
import {Link} from "react-router-dom";
import authHeaders from "../Service/request"

class AudioControlPanel extends Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
    }

    playOver(episode) {
        const name = localStorage.getItem('profile')
        fetch('/profile/playOver?name=' + name + '&episode=' + episode, authHeaders("POST"))
            .then(res => {
                $('#profiles').trigger('reload')
            })
    }

    render() {
        let source = this.props.source
        const proxy = localStorage.getItem('proxy')
        if (proxy === 'true') {
            source = '/stream' + this.props.source.replaceAll('https://', '/')
        }

        let audio = null
        if (this.props.show) {
            audio = (
                <audio controls="controls" autoPlay={this.props.play}
                    onEnded={() => this.playOver(this.props.episode)}
                ><source src={source} type="audio/x-m4a"/></audio>
            )
        }
        return audio;
    }
}

class LaterButton extends Component {
    constructor(props) {
        super(props);
    }

    listenLater(episode) {
        const name = localStorage.getItem('profile')
        fetch('/profile/playLater?name=' + name + '&episode=' + episode, authHeaders("POST"))
            .then(res => {
                $('button[action="later"][episode="' + episode + '"]').remove()

                $('#profiles').trigger('reload')
            })
    }

    render() {
        if (this.props.show) {
            return (
                <Button action="later" type="primary" size="small" episode={this.props.name}
                        onClick={() => this.listenLater(this.props.name)}>Listen Later</Button>
            );
        }
        return null
    }
}

class Episodes extends Component {
    constructor(props) {
        super(props);
        this.state = {episodes: []};
    }

    fetchEpisodes(laterPlayList) {
        const name = localStorage.getItem('profile')
        const hasProfile = (name !== "" && name != null)

        fetch('/episodes?rss=' + this.props.rss, authHeaders())
            .then(res => res.json())
            .then(res => {
                for (var i = 0; i < res.length; i++) {
                    res[i].later = hasProfile
                    for (var j = 0; j < laterPlayList.length; j++) {
                        if (laterPlayList[j].name === res[i].metadata.name) {
                            res[i].later = false
                            break
                        }
                    }
                }
                this.setState({
                    episodes: res,
                })
            })
    }

    componentDidMount() {
        const name = localStorage.getItem('profile')
        if (name === "" || !name) {
            this.fetchEpisodes([])
            return
        }
        fetch('/profiles/' + name, authHeaders())
            .then(res => res.json())
            .then(res => {
                let laterPlayList = []
                if (res.spec && res.spec.laterPlayList) {
                    laterPlayList = res.spec.laterPlayList
                }
                this.fetchEpisodes(laterPlayList)
            })
    }

    listenNow(name) {
        const episodes = this.state.episodes
        for (var i = 0; i < episodes.length; i++ ) {
            const item = episodes[i]
            if (item.metadata.name === name) {
                item.show = true
                item.play = true
            }
        }
        this.forceUpdate()
    }

    goTo(rss, name) {
        if (typeof this.props.goEpisode === 'function') {
            this.props.goEpisode(rss, name)
        }
    }

    render() {
        const {episodes} = this.state
        const rss = this.props.rss
        return (
            <div id="episodes">
                {episodes.map((item, index) => (
                    <div id={item.metadata.name} key={index} className="episode-item-in-list">
                        {item.spec.date}
                        <Link to={"/rsses/" + rss + "/episodes/" + item.metadata.name}>
                            <span className="episode-name">
                                {item.spec.title}
                            </span>
                        </Link>
                        <Button type="primary" size="small" episode={item.metadata.name}
                                onClick={() => this.listenNow(item.metadata.name)}>Listen</Button>
                        <LaterButton name={item.metadata.name} show={item.later}/>
                        <AudioControlPanel show={item.show} play={item.play} source={item.spec.audioSource} episode={item.metadata.name}/>
                    </div>
                ))}
            </div>
        )
    }
}

export default Episodes
