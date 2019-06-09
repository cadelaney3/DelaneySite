import React, { Component } from "react";
import Message from '../Message/Message';
import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';

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

class ChatHistory extends Component {
  render() {
    console.log(this.props.chatHistory);
    const { classes } = this.props;
    const messages = this.props.chatHistory.map(msg =>
      <Message message={msg.data} />
    );

    return (
      <div className="ChatHistory">
        <List className={classes.root}>
          {messages}
        </List>
      </div>
    );
  }
}

export default withStyles(styles)(ChatHistory);