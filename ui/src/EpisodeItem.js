import React, {Component} from "react";
import Button from "cuke-ui/lib/button";
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import "./EpisodeItem.css"

class EpisodeItem extends Component {
    constructor(props) {
        super(props);
        this.state = {
            object: {}
        }
    }

    goHome() {
        if (typeof this.props.goHome === 'function') {
            this.props.goHome('list')
        }
    }

    goRSS(rss) {
        if (typeof this.props.goHome === 'function') {
            this.props.goRSS(rss)
        }
    }

    componentDidMount() {
        const name = this.props.name
        fetch('/episodes/one?name=' + name)
            .then(res => res.json())
            .then(res => {
                this.setState({
                    object: res
                })
            })
    }

    render() {
        let content = ''
        if (this.state.object && this.state.object.spec) {
            content = (
                <ReactMarkdown children={this.state.object.spec.summary} remarkPlugins={[remarkGfm]} />
            )
        }
        return (
            <div className="episode-item">
                <Button type="primary" size="small" onClick={() => this.goHome()}>Go Home</Button>
                <Button type="primary" size="small" onClick={() => this.goRSS(this.props.rss)}>Go Back</Button>

                {content}
            </div>
        )
    }
}

export default EpisodeItem
