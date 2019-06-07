import React, { Component } from "react";
import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import Message from '../Message/Message';

const styles = theme => ({
  root: {
    width: '100%',
    maxWidth: 360,
    backgroundColor: theme.palette.background.paper,
  },
  inline: {
    display: 'inline',
  },
});

class Home extends Component {
    render() {
      const { classes } = this.props;
      // const messages = this.props.items;
      // console.log(messages)
      // const messages = <Message message={this.props.about} />
      const messages = this.props.items.map(msg =>
        <Message message={msg} />
      );
      return (
        <div className="Home">
          <List className={classes.root} >
            {messages}
          </List>
        </div>
      );
    }
}
  
export default withStyles(styles)(Home);