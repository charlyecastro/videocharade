import React from "react";
import { Link } from "react-router-dom";

import Constants from "./Constants";
import Alert from "./Alert.js";


export default class SignUp extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            currentUser: undefined,
            userName: "",
            email: "",
            password: "",
            confirm: "",
            errorMessage: "",
            status: 0
        };
    }

    handleSubmit(evt) {
        evt.preventDefault();

        let signUpData = {
            email: this.state.email,
            password: this.state.password,
            passwordConf: this.state.confirm,
            userName: this.state.userName,
            firstName: this.state.firstName,
            lastName: this.state.lastName
        }
        let data = JSON.stringify(signUpData)
        
         fetch("https://"+Constants.API+"/v1/users",
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
                console.log(json)
                this.props.history.push(Constants.routes.charade)
            }).catch(err => {
                console.log(err)
                return err.text().then(msg => {
                    this.setState({
                        errorMessage: msg.toString(),
                        status: err.status
                    })
                })
            })
    }

    render() {
        return (
            <div className="container">
                {this.state.errorMessage ?
                    <Alert msg={this.state.errorMessage} status={this.state.status} />
                    :
                    undefined
                }
                <header>
                </header>
                <h2> Sign Up!</h2>
                <form onSubmit={evt => this.handleSubmit(evt)}>
                    <div className="form-group">
                        <label htmlFor="name"> First Name:</label>
                        <input id="firstName" type="text" className="form-control" placeholder="Enter a your first name"
                            onInput={evt => this.setState({ firstName: evt.target.value })} />
                    </div>
                    <div className="form-group">
                        <label htmlFor="name"> Last Name:</label>
                        <input id="lastName" type="text" className="form-control" placeholder="Enter your last name"
                            onInput={evt => this.setState({ lastName: evt.target.value })} />
                    </div>
                    <div className="form-group">
                        <label htmlFor="text"> User Name:</label>
                        <input id="displayName" type="text" className="form-control" placeholder="Enter a user name"
                            onInput={evt => this.setState({ userName: evt.target.value })} />
                    </div>
                    <div className="form-group">
                        <label htmlFor="email"> Email:</label>
                        <input id="email" type="email" className="form-control" placeholder="Enter your email address"
                            onInput={evt => this.setState({ email: evt.target.value })} />
                    </div>
                    <div className="form-group">
                        <label htmlFor="password"> Password:</label>
                        <input id="password" type="password" className="form-control" placeholder="Enter your password"
                            onInput={evt => this.setState({ password: evt.target.value })} />
                    </div>
                    <div className="form-group">
                        <input id="confirm" type="password" className="form-control" placeholder="Confirm your password"
                            onInput={evt => this.setState({ confirm: evt.target.value })} />
                    </div>
                    <div className="form-group">
                        <button type="submit" className="btn"> Sign Up</button>
                    </div>
                </form>
                <p> Already have an account? <Link to={Constants.routes.signin}> Sign In </Link></p>
            </div>
        );
    }
}