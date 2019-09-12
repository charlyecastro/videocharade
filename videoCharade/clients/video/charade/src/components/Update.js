
import React from "react";

import constants from "./Constants";
import Alert from "./Alert.js";
import NavBar from "./NavBar.js";
import { Collection, CollectionItem, Row, Input, Button } from 'react-materialize'


export default class Update extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            token: window.localStorage.getItem("token"),
            user: {},
            errorMessage: "",
            update: false,
            status: 0,
            updateFirstName: "",
            updateLastName: ""
        };
    };

    handleSubmit(evt) {
        evt.preventDefault();
    }

    updateUser() {
        let updateData = {
            firstName: this.state.updateFirstName,
            lastName: this.state.updateLastName
        }
        let data = JSON.stringify(updateData)

        fetch("https://"+constants.API+"/v1/users/me",
        {
            method: "PATCH",
            body: data,
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
            this.setState({ user: json, errorMessage: "", status: 0 })
        }).catch(err => {
            return err.text().then(msg => {
                this.setState({
                    errorMessage: msg.toString(),
                    status: err.status
                })
            })
        })
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
                    console.log(err)
                    this.setState({
                        errorMessage: msg.toString(),
                        status: err.status
                    })
                })
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

        let header = {
            fontWeight: "bold",
            fontSize: "1.5rem"
        }


        return (
            <div>

                <NavBar photoURL={this.state.user.photoURL} userName={this.state.user.userName} history={this.props.history} />
                <div className="container">
                    {this.state.errorMessage ?
                        <Alert msg={this.state.errorMessage} status={this.state.status} />
                        :
                        undefined
                    }
                    <h3>Account Information</h3>

                    <Collection>
                        <CollectionItem className="row">
                            <h6 className="col offset-s1" style={header}>User Name:</h6>
                            <h6 className="col offset-s4" style={{ fontSize: "1.5rem" }}>{this.state.user.userName}</h6>
                        </CollectionItem>
                        <CollectionItem className="row">
                            <h6 className="col offset-s1" style={header}>First Name:</h6>
                            <h6 className="col offset-s4" style={{ fontSize: "1.5rem" }}>{this.state.user.firstName}</h6>
                        </CollectionItem>
                        <CollectionItem className="row">
                            <h6 className="col offset-s1" style={header}>Last Name:</h6>
                            <h6 className="col offset-s4" style={{ fontSize: "1.5rem" }}>{this.state.user.lastName}</h6>
                        </CollectionItem>
                    </Collection>

                    <Row>
                        <Input s={5} label="First Name" onInput={evt => this.setState({ updateFirstName: evt.target.value })} />
                        <Input s={5} label="Last Name" onInput={evt => this.setState({ updateLastName: evt.target.value })} />
                        <Button style={{ margin: "20px" }} onClick={() => this.updateUser()}> Update </Button>
                    </Row>

                </div>
            </div >
        );
    }
}