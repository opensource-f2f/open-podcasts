import React, {Component} from 'react';
import { useSearchParams, Link } from 'react-router-dom';
import Episodes from "./Episodes";
import Button from "cuke-ui/lib/button";
import "./RSSByCategory.css"
import Badge from "cuke-ui/lib/badge";
import authHeaders from "../Service/request"

class RSSItem extends Component {
    constructor(props) {
        super(props);
    }
    render() {
        if (this.props.image) {
            return (
                <div rss={this.props.name}>
                    <Link to={"/rsses/" + this.props.name + "/episodes"}>
                        <img width="150px" src={this.props.image} alt={this.props.name} />
                    </Link>
                    <div className="rss-title">{this.props.title}</div>
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
            categories: [],
            rssName: "",
            category: "",
        };
        this.loadEpisodes = this.loadEpisodes.bind(this);
    }

    componentDidMount() {
        const category = this.props.category
        const profile = this.props.profile
        this.loadRsses(category, profile)

        fetch('/categories', authHeaders())
            .then(res => res.json())
            .then(res => {
                if (res.items) {
                    this.setState({
                        categories: res.items,
                    })
                }
            })
    }

    loadRsses(category, profile) {
        if (!category || category === "") {
            if (!profile || profile === "") {
                fetch('/rsses', authHeaders())
                    .then(res => res.json())
                    .then(res => {
                        this.setState({
                            rsses: res,
                        })
                    })
            } else {
                fetch('/profiles/' + profile + '/subscriptions', authHeaders())
                    .then(res => res.json())
                    .then(res => {
                        this.setState({
                            rsses: res,
                        })
                    })
            }
        } else {
            fetch('/rsses?category=' + category, authHeaders())
                .then(res => res.json())
                .then(res => {
                    this.setState({
                        rsses: res,
                    })
                })
        }
    }

    shouldComponentUpdate(nextProps, nextState, nextContext) {
        if (nextProps.category !== this.props.category || nextProps.profile !== this.props.profile) {
            this.loadRsses(nextProps.category, nextProps.profile)
        }
        return true
    }

    loadEpisodes(name) {
        this.setState({rssName: name})
        if (typeof this.props.goDetail === 'function') {
            this.props.goDetail('detail', name)
        }
    }

    render() {
        const rsses = this.state.rsses
        const name = this.state.rssName
        const categories = this.state.categories

        let mySubscripted
        const profileName = localStorage.getItem('profile')
        if (profileName) {
            mySubscripted = (
                <Link to={"/rsses/subscription?profile=" + profileName}>
                    <Button type="primary" size="small">My Subscriptions</Button>
                </Link>
            )
        }

        let episodes
        if (name) {
            episodes = (
                <Episodes rss={name} key={name}/>
            )
        }

        let filter
        if (this.props.category && this.props.category !== "") {
            filter = (
                <Link to="/">
                    <Button type="primary" size="small">Go Home</Button>
                </Link>
            )
        }
        return (
            <div>
                <div>
                    {mySubscripted}
                    {filter}
                </div>
                {categories.map((item, index) => (
                    <Badge count={item.metadata.ownerReferences.length} key={index}>
                        <Link key={index} to={"/rsses/search?category=" + item.metadata.name}>
                            <Button>{item.metadata.name}</Button>
                        </Link>
                    </Badge>
                ))}
                <div id="rss_list">
                    {rsses.map((item, index) => (
                        <RSSItem name={item.metadata.name} title={item.spec.title} key={index} image={item.spec.image}/>
                    ))}
                </div>
                {episodes}
            </div>
        )
    }
}

const RSSByCategory = () => {
    const [searchParams] = useSearchParams();
    const category = searchParams.get("category")
    const profile = searchParams.get("profile")

    return (
        <RSSList category={category} profile={profile}/>
    )
}
export default RSSByCategory
