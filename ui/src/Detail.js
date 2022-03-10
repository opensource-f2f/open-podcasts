import React, {Component} from "react";
import Button from 'cuke-ui/lib/button';
import Switch from 'cuke-ui/lib/switch';
import Episodes from "./Episodes"
import "./Detail.css"
var  LocaleCode = require('locale-code')

class Detail extends Component {
    constructor(props) {
        super(props);
        this.state = {
            rsses:[],
            rssName: "",
            rss: {},
        };
        this.switchRef = React.createRef();
    }

    goTo(rss, episodeName) {
        if (rss && episodeName) {
            if (typeof this.props.goEpisode === 'function') {
                this.props.goEpisode(rss, episodeName)
            }
        } else {
            if (typeof this.props.goHome === 'function') {
                this.props.goHome('list')
            }
        }
    }

    componentDidMount() {
        fetch('/rsses')
            .then(res => res.json())
            .then(res => {
                this.setState({
                    rsses: res
                })
            })

        const name = this.props.name
        fetch('/rsses/one?name=' + name)
            .then(res => res.json())
            .then(res => {
                this.setState({
                    rss: res
                })
            })

        // load subscription
        const profile = localStorage.getItem('profile')
        if (profile === "" || !profile) {
            return
        }
        fetch('/profiles?name=' + profile)
            .then(res => res.json())
            .then(res => {
                if (res.spec && res.spec.subscription && res.spec.subscription.name) {
                    fetch('/subscriptions/one?name=' + res.spec.subscription.name)
                        .then(res => res.json())
                        .then(res => {
                            for (const rss of res.spec.rssList) {
                                if (name === rss.name) {
                                    this.switchRef.current.setState({
                                        checked: true
                                    })
                                }
                            }
                        })
                }
            })
    }

    subscribe(e) {
        const profile = localStorage.getItem('profile')
        if (profile === "" || !profile) {
            return
        }
        const requestOptions = {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' }
        };
        const rss = this.state.rss.metadata.name
        if (e) {
            fetch('/subscribe?profile=' + profile + '&rss=' + rss, requestOptions)
        } else {
            fetch('/unsubscribe?profile=' + profile + '&rss=' + rss, requestOptions)
        }
    }

    render() {
        const rss = this.props.name
        let image = ''
        if (this.state.rss.spec) {
            let language =""
            let lan = this.state.rss.spec.language
            if (lan && lan.length > 2) {
                lan = lan.replaceAll(lan.substring((lan.length - 2)), lan.substring((lan.length - 2)).toUpperCase())
                lan = LocaleCode.default.getLanguageNativeName(lan)
                language = (
                    <span>{lan}</span>
                )
            }

            image = (
                <div className="rss-head">
                    <div className="rss-head-title">
                        {this.state.rss.spec.title}
                    </div>
                    <img src={this.state.rss.spec.image} alt={this.state.rss.spec.title} width="200px"/>
                    <div className="rss-icon-group">
                        <a href={this.state.rss.spec.address} target="_blank">
                            <div className="rss-icon">
                                <div className="rss-icon-svg">
                                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512">
                                        <path d="M128.081 415.959c0 35.369-28.672 64.041-64.041 64.041S0 451.328 0 415.959s28.672-64.041 64.041-64.041 64.04 28.673 64.04 64.041zm175.66 47.25c-8.354-154.6-132.185-278.587-286.95-286.95C7.656 175.765 0 183.105 0 192.253v48.069c0 8.415 6.49 15.472 14.887 16.018 111.832 7.284 201.473 96.702 208.772 208.772.547 8.397 7.604 14.887 16.018 14.887h48.069c9.149.001 16.489-7.655 15.995-16.79zm144.249.288C439.596 229.677 251.465 40.445 16.503 32.01 7.473 31.686 0 38.981 0 48.016v48.068c0 8.625 6.835 15.645 15.453 15.999 191.179 7.839 344.627 161.316 352.465 352.465.353 8.618 7.373 15.453 15.999 15.453h48.068c9.034-.001 16.329-7.474 16.005-16.504z"></path>
                                    </svg>
                                </div>
                                <div>RSS</div>
                            </div>
                        </a>
                        {language}
                        {this.state.rss.spec.categories.map((item, index) => (
                            <span key={index}>{item}</span>
                        ))}
                        <a href={this.state.rss.spec.link} target="_blank">
                            <div className="rss-icon">
                                <div className="rss-icon-svg">
                                    <svg fill="currentColor" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
                                        <path d="M326.612 185.391c59.747 59.809 58.927 155.698.36 214.59-.11.12-.24.25-.36.37l-67.2 67.2c-59.27 59.27-155.699 59.262-214.96 0-59.27-59.26-59.27-155.7 0-214.96l37.106-37.106c9.84-9.84 26.786-3.3 27.294 10.606.648 17.722 3.826 35.527 9.69 52.721 1.986 5.822.567 12.262-3.783 16.612l-13.087 13.087c-28.026 28.026-28.905 73.66-1.155 101.96 28.024 28.579 74.086 28.749 102.325.51l67.2-67.19c28.191-28.191 28.073-73.757 0-101.83-3.701-3.694-7.429-6.564-10.341-8.569a16.037 16.037 0 0 1-6.947-12.606c-.396-10.567 3.348-21.456 11.698-29.806l21.054-21.055c5.521-5.521 14.182-6.199 20.584-1.731a152.482 152.482 0 0 1 20.522 17.197zM467.547 44.449c-59.261-59.262-155.69-59.27-214.96 0l-67.2 67.2c-.12.12-.25.25-.36.37-58.566 58.892-59.387 154.781.36 214.59a152.454 152.454 0 0 0 20.521 17.196c6.402 4.468 15.064 3.789 20.584-1.731l21.054-21.055c8.35-8.35 12.094-19.239 11.698-29.806a16.037 16.037 0 0 0-6.947-12.606c-2.912-2.005-6.64-4.875-10.341-8.569-28.073-28.073-28.191-73.639 0-101.83l67.2-67.19c28.239-28.239 74.3-28.069 102.325.51 27.75 28.3 26.872 73.934-1.155 101.96l-13.087 13.087c-4.35 4.35-5.769 10.79-3.783 16.612 5.864 17.194 9.042 34.999 9.69 52.721.509 13.906 17.454 20.446 27.294 10.606l37.106-37.106c59.271-59.259 59.271-155.699.001-214.959z"></path>
                                    </svg>
                                </div>
                            </div>
                        </a>
                        <Switch checkedChildren="取消" unCheckedChildren="收藏" ref={this.switchRef}
                                onChange={(e) => this.subscribe(e)} />
                    </div>
                    <div className="rss-head-desc">
                        {this.state.rss.spec.description}
                    </div>
                </div>
            )
        }

        return (
            <div>
                <Button type="primary" size="small" onClick={() => this.goTo()}>Go Home</Button>

                {image}
                <Episodes rss={rss} goEpisode={(rss, name) => this.goTo(rss, name)} />
            </div>
        )
    }
}

export default Detail
