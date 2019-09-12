import React from "react";

export default class Alert extends React.Component {
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
                <div className="col s12 " style = {{margin: "10px"}}>
                <div className="z-depth-2 red lighten-4" style = {{padding: "10px"}}>
                
                    <h5> Status {this.props.status}: <span>{this.props.msg}</span></h5>
                    </div>
                </div>
            </div>
        );
    }
}
