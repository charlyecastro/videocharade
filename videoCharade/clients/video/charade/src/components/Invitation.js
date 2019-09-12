import React from "react";

export default class Invitation extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
        };
    };

    render() {
        return (
            <div>
                <div className="col s12 " style = {{margin: "10px"}}>
                <div className="z-depth-2  blue lighten-3" style = {{padding: "10px"}}>
                
                    <h5> {this.props.name} wants to chat with you</h5>
                    <div className = "row">
                        <button className = "col">accept</button>
                        <button className = "col">decline</button>
                    </div>
                    </div>
                </div>
            </div>
        );
    }
}
