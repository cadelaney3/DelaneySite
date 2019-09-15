import * as React from 'react';

interface PropsType {
    children: JSX.Element;
    name: string;
    loggedIn: Float32Array;
}

export default class Master extends React.Component<PropsType, {} > {
    render() {
        console.log(this.props)
        console.log(this.props.loggedIn);
        return (
            <div>
                <button onClick={this.handleAddArticle}>
                    Add article
                </button>
            </div>
        )
    }

    handleAddArticle = () => {
        console.log("Add article");
    }
}; 

// var master = new Master();
// ReactDom.render(<Master />, document.getElementById('root'));