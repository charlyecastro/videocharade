
import React from "react";

import constants from "./Constants";
import Alert from "./Alert.js";
import Notification from "./Notification";
import NavBar from "./NavBar.js";
import ChannelBar from "./ChannelBar";
import Message from "./Message";
import CreateChannel from "./CreateChannel";


export default class Messages extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            token: window.localStorage.getItem("token"),
            user: {},
            errorMessage: "",
            status: 0,
            newMessage: "",
            currentChanID: "",
            currentChanName: "",
            currentChanDesc: "",
            chanList: [],
            messageList: [],
            createChan: false
        };
    };




    handleSubmit(evt) {
        evt.preventDefault();
        this.postMessage()
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
                console.log(json)
                this.setState({ user: json })
            }).catch(err => {
                console.log(err)
                // return err.text().then(msg => {
                //     this.setState({
                //         errorMessage: msg.toString(),
                //         status: err.status
                //     })
                // })
            })
    }

    fetchMessages(chanID) {
        console.log("fetching chan ID" + chanID)
        fetch("https://"+constants.API+"/v1/channels/" + chanID,
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
                console.log("json" + json)
                this.setState({ messageList: json.messages, })
                console.log(this.state.messageList)
                console.log(this.state.currentChanID)
            }).catch(err => {

                console.log(err)
                // return err.text().then(msg => {

                //     this.setState({
                //         messageList: [],
                //         errorMessage: msg.toString(),
                //         status: err.status
                //     })
                // })
            })
    }

    fetchChannels() {
        fetch("https://"+constants.API+"/v1/channels",
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
                console.log(json.docs)
                this.setState({
                    currentChanID: json.docs[0]._id,
                    currentChanName: json.docs[0].name,
                    currentChanDesc: json.docs[0].description,
                    chanList: json.docs,
                })
                console.log(this.state.chanList)
                console.log("current" + this.state.currentChanID)
                this.fetchMessages(this.state.currentChanID)
            }).catch(err => {
                console.log(err)
                // return err.text().then(msg => {
                //     this.setState({
                //         errorMessage: msg.toString(),
                //         status: err.status
                //     })
                // })
            })
    }

    selectChannel = (chanID, chanName, chanDesc) => {
        console.log(chanID, chanName, chanDesc)
        this.setState({
            currentChanID: chanID,
            currentChanName: chanName,
            currentChanDesc: chanDesc
        })
        this.fetchMessages(chanID)

    }

    postMessage() {
        let data = { "body": this.state.newMessage }
        data = JSON.stringify(data)
        fetch("https://"+constants.API+"/v1/channels/" + this.state.currentChanID,
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
                //console.log(json)
            }).catch(err => {
                return err.text().then(msg => {
                    this.setState({
                        errorMessage: msg.toString(),
                        status: err.status
                    })
                })
            })
    }

    componentWillMount() {
        if (!this.state.token) {
            this.props.history.push(constants.routes.signin)
        } else {
            console.log(this.state.token)
            this.fetchAccount()
            this.fetchChannels()
        }

        let sock = new WebSocket(`wss://${constants.API}/ws?auth=${this.state.token}`);
        sock.onopen = () => {
            console.log("Socket opened");
        };
        sock.onclose = () => {
            console.log("Socket closed");
        }

        sock.onmessage = (message) => {
            console.log("onmessage called!!")

            let myobj = JSON.parse(message.data)
            let eventData = myobj.data
            let eventType = myobj.type
            console.log(myobj)
            if (eventType == "message-new") {
                if (eventData.channelID == this.state.currentChanID) {
                    this.setState({ messageList: [...this.state.messageList, eventData],
                    })
                }
            } else if (eventType == "message-delete") {
                console.log("made it here")
                let myArray = this.state.messageList.filter((mess) => mess._id !== eventData)
                console.log("new array" + myArray)
                this.setState({ messageList: myArray })


            } else if (eventType == "message-update") {
                if (eventData.channelID == this.state.currentChanID) {
                    let changeIndex
                    for (let i = 0; i, this.state.messageList.length; i++) {
                        let val = this.state.messageList[i]
                        if (val._id == eventData._id) {
                            changeIndex = i
                            break
                        }
                    }

                    let updateArray = this.state.messageList;
                    updateArray[changeIndex].body = eventData.body;
                    this.setState({ messageList: updateArray })
                }
            } else if (eventType == "channel-new") {
                this.setState({ chanList: [...this.state.chanList, eventData] })
            } else if (eventType == "channel-delete") {
                console.log(eventData)
                let myArray = this.state.chanList.filter((chan) => chan._id !== eventData)
                console.log(myArray)
                this.setState({ chanList: myArray,
                }) 
                console.log(this.state.chanList)
            } else if (eventType == "channel-update") {
                console.log("made it here")
                let changeIndex
                for (let i = 0; i, this.state.chanList.length; i++) {
                    let val = this.state.chanList[i]
                    if (val._id == eventData._id) {
                        changeIndex = i
                        break
                    }
                }
                let updateArray = this.state.chanList;
                updateArray[changeIndex].name = eventData.name;
                updateArray[changeIndex].nameLower = eventData.nameLower;
                this.setState({ chanList: updateArray })
            }
       }
    }


    render() {

        return (
            <div style={{ position: "relative" }}>
                <div className="row " >
                    <NavBar photoURL={this.state.user.photoURL} userName={this.state.user.userName} history={this.props.history} />
                    <h5 style = {{paddingLeft : "10px"}}>Welcome {this.state.user.firstName}  {this.state.user.lastName}!</h5>
                    <div className="row">
                        {this.state.errorMessage ?
                            <div>
                                <a onClick={() => {
                                    this.setState({
                                        errorMessage: "",
                                        status: 0
                                    })}}> Close Banner</a>
                                <Alert className="col" msg={this.state.errorMessage} status={this.state.status} />
                            </div>
                            : undefined
                        } </div>
                    <div className="col s2">
                        {this.state.createChan ?
                            <CreateChannel userID={this.state.user.id} /> : <ChannelBar userID={this.state.user.id} chanList={this.state.chanList} onSelectChan={this.selectChannel} />
                        }
                        <a id="btn" className="waves-effect waves-light btn-large white blue-text  text-darken-4" onClick={() => this.setState({ createChan: !this.state.createChan })} >           {this.state.createChan ?
                            "Exit" : "Create Conversation"}
                        </a></div>
                    <div className="container col s10">
                       
                        <h4>Convo: {this.state.currentChanName}  </h4>

                        <div >
                            <div className="z-depth-1 " style={{ height: "60vh", paddingTop: "20px", overflowY: "auto" }}>
                                {
                                    this.state.messageList.map(mess =>
                                        <Message name={mess.creator.userName} creator={mess.creator.id} body={mess.body} key={mess._id} messID={mess._id} userID={this.state.user.id} />
                                    )
                                }
                            </div>
                        </div>
                        <form onSubmit={evt => this.handleSubmit(evt)}>
                            <div className="row center spacing">
                                <div className="input-field col s12 l10">
                                    <input type="text"
                                        className="form-control"
                                        placeholder="Message"
                                        onInput={evt => this.setState({ newMessage: evt.target.value })} />
                                </div>
                                <div className="form-group col s2" style={{ marginTop: "10px" }}>
                                    <button type="submit" className="btn" > Send </button>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        );
    }
}


