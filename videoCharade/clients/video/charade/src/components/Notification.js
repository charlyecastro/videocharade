import React from "react";

export default class Notification extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            email: "",
            password: ""
        };
    };

    render() {
        return (
            <div>
                <div className="col s12" style = {{margin: "10px"}}>
                <div className="z-depth-2  teal accent-1" style = {{padding: "10px"}}>
                    <h5> Notification: {this.props.msg}</h5>
                    </div>
                </div>
            </div>
        );
    }
}
