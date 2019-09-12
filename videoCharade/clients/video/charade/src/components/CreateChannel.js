import React from "react";

import constants from "./Constants";
import UserSearch from "./UserSearch";
import { Input } from 'react-materialize'

export default class CreateChannel extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            token: window.localStorage.getItem("token"),
            chanName: "",
            chanDesc: "",
            chanPriv: false,
            idList: [],
            addedUser : false
        };
    };


    handleSubmit(evt) {
        evt.preventDefault();
        console.log("target" + evt.target.value)

        let mySet = new Set(this.state.idList)
        console.log(mySet)
        let uniqueArray = Array.from(mySet);
        console.log(uniqueArray)

        let channelData = {
            name: this.state.chanName,
            description: this.state.chanDesc,
            members:  uniqueArray,
            private: this.state.chanPriv
        }
        let data = JSON.stringify(channelData)

        console.log("data being sent: " + data)

         fetch("https://"+constants.API+"/v1/channels",
            {
                method: "POST",
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
                console.log(json)
            }).catch(err => {
                console.log(err)
                return err.text().then(msg => {
                    this.setState({ errorMessage: msg.toString(),
                        status: err.status})
                })
            })
    }

    componentWillMount() {
        this.setState({ chanPriv: false })
    }

    addUser = (userId) => {
        if (userId != this.props.userID) {
        this.state.idList.push({id : userId})
        let mySet = new Set(this.state.idList)
        let uniqueArray = Array.from(mySet);
        console.log(uniqueArray)
        console.log("list after: " + this.state.idList)
        }
        this.setState({addedUser: true})
    }

    render() {

        return (
            <div className="z-depth-1 " style={{ backgroundColor: "white", padding: "20px" }}>

                    {this.state.addedUser ?
                            <div className = "row light-blue accent-1">
                            <div className = "col"> Added User!</div>
                            <a  style = {{color : "black"}}className = "col" onClick={() => {
                                    this.setState({
                                        addedUser: false,
                                    })}}> close </a>
                            </div>
                            : undefined
                        } 
                <h5>Create a Conversation</h5>

                <form className="container" style={{ paddingTop: "10px" }} onSubmit={evt => this.handleSubmit(evt)}>

                    {/* <div style={{ width: "100%" }} className="form-group">
                        <Input type='checkbox' label='Private' onChange={() => this.setState({ chanPriv: !this.state.chanPriv })} />
                    </div>
                    {this.state.chanPriv ?
                        <UserSearch token={this.state.token} onAddUser={this.addUser} /> : undefined
                    } */}
                    <div style={{ paddingTop: "10px" }} className="form-group">
                        <label htmlFor="text"> Name:</label>
                        <input id="name" type="text" name = "chanName"
                            className="form-control" placeholder="enter a new name"
                            onInput={evt => this.setState({ chanName: evt.target.value })}
                        />
                    </div>
                    <div style={{ paddingTop: "10px" }} className="form-group">
                        <label htmlFor="text">Conversation desciption:</label>
                        <input name = "chanDesc" id="descripton" type="text" className="form-control" placeholder="enter a new description"
                            onInput={evt => this.setState({ chanDesc: evt.target.value })}
                        />
                    </div>
                    <div className="form-group">
                        <button type="submit" className="btn" > Create </button>
                    </div>
                </form>


            </div>
        );
    }
}