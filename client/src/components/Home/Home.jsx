import React, { Component } from "react";
import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/List';
import Grid from '@material-ui/core/Grid';
import ButtonBase from '@material-ui/core/ButtonBase';
import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import mugshot from '../../images/mugshot.jpg';

const styles = theme => ({
  root: {
    flexGrow: 1,
    backgroundColor: theme.palette.background.paper,
    paddingTop: theme.spacing.unit * 8,
    //paddingLeft: theme.spacing.unit * 4,
  },
  image: {
      width: 256,
      height: 256,
  },
  img: {
      margin: 'auto',
      display: 'block',
      maxWidth: '100%',
      maxHeight: '100%',
  },
  firstColPaper: {
      width: 256,
      padding: 8,
      margin: 'auto',
  },
  contentPaper: {
      padding: 8,
      margin: 'auto',
  },
  inline: {
    display: 'inline',
  },
});

class Home extends Component {
    render() {
      const { classes } = this.props;   
      console.log(this.props);
      console.log(this.props.items.body);
      const facts = this.props.items.facts.map(fact =>
        <ListItem key={fact}>
            <Paper className={classes.contentPaper}>
                <Typography variant="body2" align="left">
                    {fact}
                </Typography>
            </Paper>
        </ListItem>
      );
      return (
            <Grid container className={classes.root} spacing={8}>
                <Grid item xs={3}>
                    <Grid container direction="column" spacing={8}>
                        <Grid item>
                            <ButtonBase className={classes.image}>
                                <img className={classes.img} alt="Chris" src={mugshot} />
                            </ButtonBase>
                        </Grid>
                        <Grid item>
                            <Paper className={classes.firstColPaper}>
                                <Typography variant="body2" align="left">
                                    github.com/cadelaney3
                                </Typography>
                                <Typography variant="body2" align="left">
                                    linked.com/in/cadelaney3
                                </Typography>
                            </Paper>
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item xs={12} sm container>
                    <Grid item xs container direction="column" spacing={8}>
                        <List>
                            <ListItem key={this.props.items.body}>
                                <Paper className={classes.contentPaper}>
                                    <Typography variant="body1" align="left">
                                        {this.props.items.body}
                                    </Typography>
                                </Paper>
                            </ListItem>
                            {facts}
                        </List>
                    </Grid>
                </Grid>
            </Grid>
      );
    }
}
  
export default withStyles(styles)(Home);