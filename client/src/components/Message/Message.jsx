import React, { Component } from 'react';
import { withStyles } from '@material-ui/core/styles';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import Avatar from '@material-ui/core/Avatar';

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

class Message extends Component {
    constructor(props) {
        super(props);
        console.log(this.props.message)
        let temp = JSON.parse(JSON.stringify(this.props.message));
        console.log("This is temp: ", temp.body)
        this.state = {
            message: temp
        };
    }

    render() {
        //const { classes } = this.props;
        return (
            <ListItem alignItems="flex-start">
                <ListItemAvatar>
                    <Avatar alt="Remy Sharp" src="/static/images/avatar/1.jpg" />
                </ListItemAvatar>
                <ListItemText
                    primary={this.state.message.body}
                />              
            </ListItem>
        ) 
    }
}

export default withStyles(styles)(Message);