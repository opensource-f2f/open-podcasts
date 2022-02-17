import React, { Component } from 'react';

class Profile extends Component {
    constructor(props) {
        super(props);
        this.state = { laterPlayList: [] };
    }

    componentDidMount() {
        fetch('/profiles?name=linuxsuren')
            .then(res => res.json())
            .then(res => {
                this.setState({
                    laterPlayList: res.spec.laterPlayList
                })
            })
    }

    play() {
        console.log('play')
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
            </div>
        );
    }
}

export default Profile;
