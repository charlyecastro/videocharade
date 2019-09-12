
import React from "react";

import constants from "./Constants";
import Alert from "./Alert.js";
import NavBar from "./NavBar.js";
import SummmaryCard from "./SummaryCard";

export default class Summary extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            token: window.localStorage.getItem("token"),
            user: {},
            errorMessage: "",
            status: 0,
            query: "",
        };
    };

    handleSubmit(evt) {
        evt.preventDefault();
    }

    fetchAccount() {
         fetch("https://"+constants.API+"/v1/users/me",
            {
                method: "GET",
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': this.state.token,
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                return res.json();
            }).then(json => {
                this.setState({ user: json })
            }).catch(err => {
                return err.text().then(msg => {
                    this.setState({
                        errorMessage: msg.toString(),
                        status: err.status
                    })
                }).catch(err => {
                })
            })
    }

    fetchSummary() {
        fetch("https://"+constants.API+"/v1/summary?url=" + this.state.query, {
            method: "GET",
            headers: {
                "Accept": "application/json"
            }
        }).then(res => {
            if (res.status > 200) {
                throw res;
            }
            return res.json();
        }).then(json => {
            this.setState({ summary: json })
        }).catch(err => {
            return err.text().then(msg => {
                this.setState({
                    errorMessage: msg.toString(),
                    status: err.status
                })
            }
            )
        })
    }

    componentDidMount() {
        if (!this.state.token) {
            this.props.history.push(constants.routes.signin)
        } else {
            this.fetchAccount()
        }

    }

    render() {

        return (
            <div>
                <NavBar photoURL={this.state.user.photoURL} userName={this.state.user.userName} history={this.props.history} />
                <div className="container">
                    
                    {this.state.errorMessage ?
                        <Alert msg={this.state.errorMessage} status={this.state.status} />
                        :
                        undefined
                    }
                    <div className="row center spacing">
                        <div className="input-field col s12  l9">
                            <input placeholder="http://url.com" id="input" type="text" className="validate" onInput={evt => this.setState({ query: evt.target.value })} />
                            </div>
                            <div className="input-field col s2 offset-s3  l2">
                                <a id="btn" className="waves-effect waves-light btn-large white blue-text  text-darken-4" onClick={() => this.fetchSummary()} >Summarize</a>

                            </div>
                        
                    </div>
                    <div className="divider"></div>
                    <h5 className="blue-text  text-darken-4 ">result:</h5>
                    <div id="result" className="container">
                    </div>
                    {this.state.summary ?
                        <SummmaryCard summary={this.state.summary} />
                        :
                        undefined
                    }
                </div>
            </div>

        );
    }
}