import './GlobalAudio.css'
import React, { Component } from "react";
import ReactDOM from 'react-dom'
import ReactJkMusicPlayer from 'react-jinke-music-player'
import 'react-jinke-music-player/assets/index.css'
import $ from "jquery";

class GlobalAudio extends Component {
    constructor(props) {
        super(props);
        this.state = {
            mode: "native1",
            playList: [],
            clearPriorAudioLists: false,
        }
        this.useMode = this.useMode.bind(this);
        this.play = this.play.bind(this);
        this.loadPlayList = this.loadPlayList.bind(this);
        this.audioInstance = null
    }

    useMode(mode) {
        this.setState({
            mode: mode,
        })
    }

    play(src) {
        if (this.state.mode === "native") {
            return this.playInNative(src)
        } else {
            this.audioInstance.play()
            return null
        }
    }

    playInNative(src) {
        const parent = $('#global-audio-zone').show()
        const source = $(document.createElement('source'))
        source.attr('src', src)
        source.attr('type', 'audio/x-m4a')

        const audio = $(document.createElement('audio'))
        audio.attr('controls','')
        audio.append(source)
        parent.empty().append(audio)
        audio.trigger('play')
        return audio
    }

    onAudioEnded(currentPlayId,audioLists,audioInfo) {
        const id = audioInfo.id

        const name = localStorage.getItem('profile')
        if (name === "" || !name) {
            return
        }
        $.post('/profile/playOver?name=' + name + '&episode=' + id)
    }

    componentDidMount() {
        const name = localStorage.getItem('profile')
        if (name === "" || !name) {
            return
        }
        const comObject = this
        fetch('/profiles?name=' + name)
            .then(res => res.json())
            .then(res => {
                if (res.spec && res.spec.laterPlayList) {
                    comObject.loadPlayList(res.spec.laterPlayList)
                }
            })
    }

    loadPlayList(laterPlayList) {
        if (laterPlayList.length === 0) {
            return
        }

        async function fetchLaterPlayList() {
            const audioList = []
            for (let i = 0; i < laterPlayList.length; i++) {
                const episode =  laterPlayList[i]

                await $.getJSON('/episodes/one?name=' + episode.name, function (item) {
                    let source = item.spec.audioSource
                    const proxy = localStorage.getItem('proxy')
                    if (proxy === 'true') {
                        source = '/stream' + source.replaceAll('https://', '/')
                    }

                    let author
                    if (item.metadata.annotations && item.metadata.annotations["title.podcast"]) {
                        author = item.metadata.annotations["title.podcast"]
                    }

                    audioList.push({
                        id: item.metadata.name,
                        singer: author,
                        name: item.spec.title,
                        musicSrc: source,
                        cover: item.spec.coverImage,
                    })
                })
            }
            return audioList
        }
        const comObject = this;
        fetchLaterPlayList().then((audioList) => {
            comObject.setState({
                playList: audioList,
            })
        })
    }

    render() {
        if (this.state.mode === "native") {
            return (
                <div id="global-audio-zone" className="global-audio"></div>
            )
        } else {
            const clearPriorAudioLists = true
            const theme = 'light'
            const params = {
                clearPriorAudioLists: clearPriorAudioLists,
                theme: theme,
                preload: false,
                onAudioEnded: this.onAudioEnded,
                audioLists: this.state.playList,
                remember: true,
                autoPlay: false,
            }
            return (
                <ReactJkMusicPlayer mode="full"
                                    {...params}
                    getAudioInstance={(instance) => {
                        this.audioInstance = instance
                    }}
                />
            )
        }
    }
}

let globalAudio = document.getElementById('globalAudio')
if (!globalAudio) {
    globalAudio = document.createElement('div')
    globalAudio.setAttribute('id', 'globalAudio')
    document.body.appendChild(globalAudio)
}
let globalAudioCom = ReactDOM.render(React.createElement(
    GlobalAudio
), globalAudio)

export default globalAudioCom
