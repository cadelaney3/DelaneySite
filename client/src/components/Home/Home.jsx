import React, { Component } from "react";
import { withStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/List';
import Grid from '@material-ui/core/Grid';
import ButtonBase from '@material-ui/core/ButtonBase';
import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import Link from '@material-ui/core/Link';
import Octicon, {MarkGithub, File} from '@primer/octicons-react';
import Avatar from '@material-ui/core/Avatar';
import mugshot from '../../images/mugshot.jpg';
import linkedinLogo from '../../images/In-Blue-26@2x.png';
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import Resume from '../../files/ChrisDelaney_Resume.pdf';

const styles = theme => ({
  root: {
    flexGrow: 1,
    paddingTop: theme.spacing(4),
  },
  image: {
      width: 270,
      height: 270,
  },
  img: {
      margin: 'auto',
      display: 'block',
      maxWidth: '100%',
      maxHeight: '100%',
  },
  firstColPaper: {
      width: 270,
      padding: 4,
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
    align: 'left'
  },
  avatar: {
      width: "25px",
      height: "25px",
      marginTop: '5px',
      marginRight: '8px'
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
      marginRight: '8px',
      verticalAlign: 'center',
      color: 'white',
      marginTop: '5px'
  },
  fab: {
      margin: theme.spacing(1),
      position: 'fixed',
  }
});

class Home extends Component {
    constructor(props) {
        super(props);

        this.state = {
            loggedIn: sessionStorage.getItem("loggedIn"),
        }
    }
    
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
                        <Paper className={classes.firstColPaper} align='left'>
                                <Link href="https://github.com/cadelaney3/" className={classes.link} target="_blank" rel="noopener" >
                                    <Octicon className={classes.octicon} icon={MarkGithub} ariaLabel="GitHub" noWrap />
                                    {"https://github.com/cadelaney3"}
                                </Link>
                                <Link href="https://linked.com/in/cadelaney3/" className={classes.link} target="_blank" rel="noopener">
                                <Grid container alignItems="center">
                                    <Avatar className={classes.avatar} alt="in" src={linkedinLogo} inline="true" />
                                    {"https://linkedin.com/in/cadelaney3"}
                                </Grid>
                                </Link>
                                <Link href={Resume} className={classes.link} target="_blank" rel="noopener">
                                    <Octicon className={classes.octicon} icon={File} ariaLabel="File" noWrap />
                                    {"My Resume"}
                                </Link>
                        </Paper>
                    </Grid>
                </Grid>
                <Grid item>
                    <Paper className={classes.firstColPaper} align='left'>
                            <Link href="https://github.com/cadelaney3/" className={classes.link} target="_blank" rel="noopener" >
                                <Octicon className={classes.octicon} icon={MarkGithub} ariaLabel="GitHub" noWrap />
                                {"https://github.com/cadelaney3"}
                            </Link>
                            <Link href="https://linked.com/in/cadelaney3/" className={classes.link} target="_blank" rel="noopener">
                            <Grid container alignItems="center">
                                <Avatar className={classes.avatar} alt="in" src={linkedinLogo} inline="true" />
                                {"https://linkedin.com/in/cadelaney3"}
                            </Grid>
                            </Link>
                    </Paper>
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
            <Fab aria-label="Add" className={classes.fab}>
                <AddIcon />
            </Fab>
        </Grid>
      );
    }
}
  
export default withStyles(styles)(Home);