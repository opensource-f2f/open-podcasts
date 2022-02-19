import React, { Component } from 'react'
import './Episodes.css'
import $ from 'jquery'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'

class AudioControlPanel extends Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
    }

    playOver(episode) {
        const requestOptions = {
            method: 'POST',
        };
        const name = localStorage.getItem('profile')
        fetch('/profile/playOver?name=' + name + '&episode=' + episode, requestOptions)
            .then(res => {
                $('#profiles').trigger('reload')
            })
    }

    render() {
        let audio = null
        if (this.props.show) {
            audio = (
                <audio controls="controls" autoPlay={this.props.play}
                    onEnded={() => this.playOver(this.props.episode)}
                ><source src={this.props.source} type="audio/x-m4a"/></audio>
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
        const requestOptions = {
            method: 'POST',
        };
        const name = localStorage.getItem('profile')
        fetch('/profile/playLater?name=' + name + '&episode=' + episode, requestOptions)
            .then(res => {
                $('button[action="later"][episode="' + episode + '"]').remove()

                $('#profiles').trigger('reload')
            })
    }

    render() {
        if (this.props.show) {
            return (
                <button action="later" episode={this.props.name} onClick={() => this.listenLater(this.props.name)}>Listen Later</button>
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

        fetch('/episodes?rss=' + this.props.rss)
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
        if (name === "" || !name || name == null) {
            this.fetchEpisodes([])
            return
        }
        fetch('/profiles?name=' + name)
            .then(res => res.json())
            .then(res => {
                const laterPlayList = res.spec.laterPlayList
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

    render() {
        const {episodes} = this.state
        return (
            <div id="episodes">
                {episodes.map((item, index) => (
                    <div id={item.metadata.name} src="" key={index}>
                        <span>{item.spec.title}</span><br/>
                        <ReactMarkdown children={item.spec.summary} remarkPlugins={[remarkGfm]} />
                        <button episode={item.metadata.name} onClick={() => this.listenNow(item.metadata.name)}>Listen</button>
                        <LaterButton name={item.metadata.name} show={item.later}/>
                        <AudioControlPanel show={item.show} play={item.play} source={item.spec.audioSource} episode={item.metadata.name}/>
                    </div>
                ))}
            </div>
        )
    }
}

export default Episodes
