import React, { Component } from 'react';
import $ from 'jquery'
import './Profile.css'
import avatar from './images/img_avatar.png'
import Modal from 'react-modal'

function createAudio(src, parent) {
    const source = $(document.createElement('source'))
    source.attr('src', src)
    source.attr('type', 'audio/x-m4a')

    const audio = $(document.createElement('audio'))
    audio.attr('controls','')
    audio.append(source)
    parent.append(audio)
    return audio
}

function playEpisode(episdoe, callback) {
    $.getJSON('/episodes/one?name=' + episdoe, function (item){
        createAudio(item.spec.audioSource, $('#global-audio-zone').show()).trigger('play').on('ended', function () {
            const episode = $(this).attr('episode')
            const profile = window.localStorage.getItem('profile')
            $.post('/profile/playOver?name=' + profile + '&episode=' + episode, function (){
                $('span[episode=' + episode + ']').remove()

                if (callback) {
                    callback()
                }
            })
        }).attr('episode', item.metadata.name)
    })
}

Modal.setAppElement('#root');
class ProfileModal extends Component {
    constructor(props) {
        super(props);
        this.state = {
            isOpen: false,
            rssURL: "",
            rssName: ""
        }
    }

    closeModal() {
        this.setState({isOpen: false})
    }

    addRSS() {
        console.log(this.state)
        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                address: this.state.rssURL,
                name: this.state.rssName
            })
        };
        fetch('/rsses', requestOptions)
    }
    setRSSURL(e) {
        this.setState({
            rssURL: e.target.value
        })
    }
    setRSSName(e) {
        this.setState({
            rssName: e.target.value
        })
    }

    render() {
        return (
            <div>
                <Modal
                    isOpen={this.state.isOpen}
                    contentLabel="Example Modal"
                >
                    <button onClick={() => this.closeModal()}>Close</button>
                    <div>
                        New RSS feed:<input onChange={(e) => this.setRSSURL(e)}/> with name:
                        <input onChange={(e) => this.setRSSName(e)}/><button onClick={() => this.addRSS()}>Add</button>
                    </div>
                </Modal>
            </div>
        );
    }
}

class Profile extends Component {
    constructor(props) {
        super(props);
        this.state = {
            laterPlayList: [],
            showModal: false,
        };
        this.profileModalElement = React.createRef();
    }

    reload(){
        const name = localStorage.getItem('profile')
        if (name === "" || !name || name == null) {
            return
        }
        fetch('/profiles?name=' + name)
            .then(res => res.json())
            .then(res => {
                this.setState({
                    laterPlayList: res.spec.laterPlayList
                })
            })
    }

    componentDidMount() {
        this.reload()
        const com = this
        $('#profiles').on('reload', function (){
            com.reload()
        })
    }

    play() {
        $('span[episode]').each(function (i, item){
            item = $(item)
            const episode = item.attr('episode')
            $('#' + episode).trigger('play-audio').on('ended-audio', function () {
                item.remove()
                // playButton.click()
            })

            // hideGlobalAudioZone()
            if ($('#' + episode).length > 0) {
                $([document.documentElement, document.body]).animate({
                    scrollTop: $('#' + episode).offset().top
                }, 2000);
            } else {
                playEpisode(episode, function () {
                })
            }
            return false
        })
    }

    toggleModal() {
        this.profileModalElement.current.setState({isOpen: true})
    }

    render() {
        const {laterPlayList} = this.state;
        return (
            <div id="profiles">
                <div style={{display: "none"}} id="login-zone">
                    <label>
                        Name: <input name="name" id="login-name" />
                    </label>
                    <div><button action="login">Login</button></div>
                    <div><button action="register">Register</button></div>
                </div>
                <div>
                    <div>
                        <span>Listen Later List: </span>
                        {laterPlayList.map((item, index) => (
                            <span episode={item.name} key={index}>{item.name}</span>
                            )
                        )}
                        <button onClick={this.play}>Play</button>
                    </div>
                </div>

                <img src={avatar} className="avatar" alt="Avatar" onClick={() => this.toggleModal()}/>

                <ProfileModal ref={this.profileModalElement}/>
            </div>
        );
    }
}

export default Profile;
