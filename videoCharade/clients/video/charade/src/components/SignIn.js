import React from "react";
import { Link } from "react-router-dom";

import Constants from "./Constants";
import Alert from "./Alert.js";

export default class SignIn extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            //currentUser: undefined,
            email: "",
            password: "",
            errorMessage: "",
            status: 0
        };
    };

    componentDidMount(){
        let token = window.localStorage.getItem("token")
        if (token != null) {
            this.props.history.push(Constants.routes.charade)
        }
        // console.log("Home has access to bearer: " + token)
        // this.setState({ user: this.props.location.state })
    }


    handleSubmit(evt) {
        evt.preventDefault();

        let signInData = {
            email: this.state.email,
            password: this.state.password,
        }

        let data = JSON.stringify(signInData)

         fetch("https://"+Constants.API+"/v1/sessions",
            {
                method: "POST",
                body: data,
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                })
            }).then(res => {
                if (res.status > 300) {
                    throw res;
                }
                let token = res.headers.get("Authorization");
                window.localStorage.setItem("token", token);
                return res.json();
            }).then(json => {
                this.props.history.push(Constants.routes.charade)
            }).catch(err => {
                console.log(err)
                return err.text().then(msg => {
                    this.setState({ errorMessage: msg.toString(),
                        status: err.status})
                })
            })
    }

    render() {
        return (
            <div className="container">
                {this.state.errorMessage ?
                        <Alert msg = {this.state.errorMessage}  status = {this.state.status}/>
                        :
                    undefined
                }
                <header>
                </header>
                <h2> Sign In </h2>

                <form onSubmit={evt => this.handleSubmit(evt)}>
                    <div className="form-group">
                        <label htmlFor="email">Email:</label>
                        <input id="email" type="email"
                            className="form-control" placeholder="enter your email address"
                            onInput={evt => this.setState({ email: evt.target.value })}
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="password">Password:</label>
                        <input id="password" type="password" className="form-control" placeholder="enter your password"
                            onInput={evt => this.setState({ password: evt.target.value })}
                        />
                    </div>
                    <div className="form-group">
                        <button type="submit" className="btn" > Sign In
    </button>
                    </div>
                </form>

                <p> Dont have an account? <Link to={Constants.routes.signup}> Create One! </Link> </p>

            </div>

        );
    }
}
