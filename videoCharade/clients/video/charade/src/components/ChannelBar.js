import React from "react";

import constants from "./Constants";

export default class ChannelBar extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            selectChan: "",

        };
    };

    handleSelected = (chan, name, desc) => {
        console.log("channel Name: " + chan)
        this.props.onSelectChan(chan, name, desc);
    }

    editChannel(chanID) {

        let newChannel = prompt("Please enter your new channel")
        console.log(newChannel)
        let data = { "name": newChannel }
        console.log(data)
        data = JSON.stringify(data)

        if (newChannel != null) {

              fetch("https://"+constants.API + "/v1/channels/" + chanID,
                {
                    method: "PATCH",
                    body: data,
                    headers: new Headers({
                        'Content-Type': 'application/json',
                        'Accept': 'application/json',
                        'Authorization': window.localStorage.getItem("token"),
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
                })
        }
    }

    deleteChannel(chanID) {
         fetch("https://"+ constants.API + "/v1/channels/" + chanID,
            {
                method: "DELETE",
                headers: new Headers({
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    'Authorization': window.localStorage.getItem("token"),
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
            })
    }


    render() {
        return (
            <div className="z-depth-1" style={{ padding: "10px" }}>


                <h5> Conversations</h5>
                {
                    this.props.chanList.map(channel =>
                        // <h6> {channel.name}</h6>

                        <div key={channel._id} className="row " style={{ backgroundColor: "#26a69a", padding: "5px", borderRadius: "2px" }} >

                            <a style={{ color: "#FFFFFF", fontWeight: "bold" }} onClick={() => this.handleSelected(channel._id, channel.name, channel.description)} className="col" >
                            {channel.private ?
                                <span className="col" style={{ fontWeight: "bold", color: "white", padding : "0px" }}>*</span>
                                : undefined}
                            {channel.name} </a>

                            {this.props.userID == channel.creator.id ?
                                <div>
                                    <a style={{ color: "#FFFFFF" }} onClick={() => this.editChannel(channel._id)} className="col" >edit </a>
                                    <a style={{ color: "#FFFFFF" }} onClick={() => this.deleteChannel(channel._id)} className="col" >delete </a>
                                </div>
                                : undefined
                            }


                        </div>
                    )
                }
            </div>
        );
    }
}