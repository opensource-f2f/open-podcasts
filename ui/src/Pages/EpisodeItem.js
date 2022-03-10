import React, {Component} from "react";
import Button from "cuke-ui/lib/button";
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import "./EpisodeItem.css"
import {Link, useParams} from "react-router-dom";

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
        let coverImage = ''
        if (this.state.object && this.state.object.spec) {
            content = (
                <ReactMarkdown children={this.state.object.spec.summary} remarkPlugins={[remarkGfm]} />
            )
            if (this.state.object.spec.coverImage) {
                coverImage = (
                    <img src={this.state.object.spec.coverImage} alt={this.state.object.metadata.name} width="200"/>
                )
            }
        }
        return (
            <div className="episode-item">
                <Link to="/">
                    <Button type="primary" size="small">Go Home</Button>
                </Link>
                <Link to={"/rsses/" + this.props.rss + "/episodes"}>
                    <Button type="primary" size="small">Go Back</Button>
                </Link>

                <div>
                    {coverImage}
                    {content}
                </div>
            </div>
        )
    }
}

export default function () {
    let { rss, episode } = useParams();
    return <EpisodeItem name={episode} rss={rss} />;
}
