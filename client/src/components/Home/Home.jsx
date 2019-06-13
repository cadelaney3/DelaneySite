import React, { Component } from "react";
import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/List';
import Grid from '@material-ui/core/Grid';
import ButtonBase from '@material-ui/core/ButtonBase';
import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import Link from '@material-ui/core/Link';
import Octicon, {MarkGithub} from '@primer/octicons-react';
import Avatar from '@material-ui/core/Avatar';
import mugshot from '../../images/mugshot.jpg';
import linkedinLogo from '../../images/In-Blue-26@2x.png';

const styles = theme => ({
  root: {
    flexGrow: 1,
    paddingTop: theme.spacing.unit * 4,
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
  link: {
    textDecoration: 'none',
    color: '#ffd600',
  },
  avatar: {
      width: "25px",
      height: "25px",
      marginTop: '5px',
  },
  typography: {
      paddingTop: "8px",
      paddingBottom: "8px",
      color: "#ffd600",
      fontFamily: 'Roboto',
  },
  body: {
      fontFamily: 'sans-serif',
      color: '#fff176',
  },
  octicon: {
      height: '25px',
      width: '25px',
      paddingX: '3px',
      verticalAlign: 'center',
      color: 'white',
  },
});

class Home extends Component {
    render() {
      const { classes } = this.props;   
      const facts = this.props.items.facts.map(fact =>
        <ListItem key={fact}>
            <Paper className={classes.contentPaper}>
                <Typography className={classes.body} variant="body2" align="left">
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
                                    <Link href="https://github.com/cadelaney3/" className={classes.link} target="_blank" rel="noopener" >
                                        <Octicon className={classes.octicon} icon={MarkGithub} ariaLabel="GitHub" noWrap />
                                        {"https://github.com/cadelaney3"}
                                    </Link>
                                    <Link href="https://linked.com/in/cadelaney3/" className={classes.link} target="_blank" rel="noopener">
                                    <Grid container justify="center" alignItems="center">
                                        <Avatar className={classes.avatar} alt="in" src={linkedinLogo} inline="true" />
                                        {"https://linkedin.com/in/cadelaney3"}
                                    </Grid>
                                    </Link>
                            </Paper>
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item xs={12} sm container>
                    <Grid item xs container direction="column" spacing={8}>
                        <Typography className={classes.typography} variant="h2" align="left">
                            Bio
                        </Typography>
                        <List>
                            <ListItem key={this.props.items.body}>
                                <Typography className={classes.body} variant="body1" align="left">
                                    {this.props.items.body}
                                </Typography>
                            </ListItem>
                        </List>
                        <Typography className={classes.typography} variant="h2" align="left">
                            Fun Facts
                        </Typography>
                        <List>
                            {facts}
                        </List>                      
                    </Grid>
                </Grid>
            </Grid>
      );
    }
}
  
export default withStyles(styles)(Home);