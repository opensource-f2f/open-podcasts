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
            let lan = this.state.rss.spec.language
            if (lan.length > 2) {
                lan = lan.replaceAll(lan.substring((lan.length - 2)), lan.substring((lan.length - 2)).toUpperCase())
            }
            const language = LocaleCode.default.getLanguageNativeName(lan)

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
                        <span>{language}</span>
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
