import React, { useState, useEffect } from 'react';
// import clsx from 'clsx';
import { makeStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
// import CssBaseline from '@material-ui/core/CssBaseline';
import List from '@material-ui/core/List';
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import InboxIcon from '@material-ui/icons/MoveToInbox';
import MailIcon from '@material-ui/icons/Mail';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
// import CardMedia from '@material-ui/core/CardMedia';
import Button from '@material-ui/core/Button';
import AddArticle from '../AddArticle/AddArticle';

const drawerWidth = 240;

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
  content: {
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'normal',
    // flexGrow: 1,
    padding: theme.spacing(3),
  },
  card: {
    width: 450,
    margin: theme.spacing(1)
  },
  toolbar: theme.mixins.toolbar,
}));

function ArticleCards(props) {
  const classes = useStyles();
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [articles, setArticles] = useState([]);

  const getResults = () => {
    //fetch("http://localhost:8080/home")
    fetch("http://172.26.34.14:8080/articles")
    .then(res => res.json())
    .then(result => {
        setArticles(result);
        if(articles !== []) {
          setIsLoaded(true);
        }
    })
    .catch(error => {
        setIsLoaded(true);
        setError(error);
    })
  };

  useEffect(() => {
    getResults();
  }); //, [newArticle]);

  useEffect(() => {
    if (articles.article) {
      console.log(articles);
      props.setFeed(articles.article.map(item =>
        <Card className={classes.card} key={item.title}>
          <CardActionArea>
            <CardContent align="left">
              <Typography gutterBottom variant="h5" component="h2">
                {item.title}
              </Typography>
              <Typography variant="body2" color="textSecondary" component="p">
                {item.description}
              </Typography>
            </CardContent>
          </CardActionArea>
          <CardActions>
            <Button size="small" color="primary">
              Share
            </Button>
            <Button size="small" color="primary">
              Learn More
            </Button>
          </CardActions>
        </Card>
      ));
    } else {
      props.setFeed(
        <div>
          Loading...
        </div>
      )
    }
  }, [isLoaded]);

  //return (
  if (error) {
    return (
      <div>
        Error: {error.message}
      </div>
    );
  } else if (!isLoaded) {
    return (
      <div>
        Loading...
      </div>
    );
  } else {
    return props.feed
  }
}

export default function Articles(props) {
  const classes = useStyles();
  const [newArticle, setNewArticle] = useState(false);
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [articles, setArticles] = useState([]);
  const [feed, setFeed] = useState(null);

  const getResults = () => {
    //fetch("http://localhost:8080/home")
    fetch("http://172.26.34.14:8080/articles")
    .then(res => res.json())
    .then(result => {
        setArticles(result);
       // if( {articles} !== []) {
        setIsLoaded(true);
       // }
    })
    .catch(error => {
        setIsLoaded(true);
        setError(error);
    })
  };

  useEffect(() => {
    getResults();
  }, [newArticle]);

  useEffect(() => {
    if (articles.article) {
      console.log(articles);
      setFeed(articles.article.map(item =>
        <Card className={classes.card} key={item.title}>
          <CardActionArea>
            <CardContent align="left">
              <Typography gutterBottom variant="h5" component="h2">
                {item.title}
              </Typography>
              <Typography variant="body2" color="textSecondary" component="p">
                {item.description}
              </Typography>
            </CardContent>
          </CardActionArea>
          <CardActions>
            <Button size="small" color="primary">
              Share
            </Button>
            <Button size="small" color="primary">
              Learn More
            </Button>
          </CardActions>
        </Card>
      ));
    }    
  }, [articles]);

  var page =
    <div className={classes.root} align='center'>
      <Drawer
        className={classes.drawer}
        variant="permanent"
        classes={{
          paper: classes.drawerPaper,
        }}
      >
        <div className={classes.toolbar} />
        <List>
          {['Inbox', 'Starred', 'Send email', 'Drafts'].map((text, index) => (
            <ListItem button key={text}>
              <ListItemIcon>{index % 2 === 0 ? <InboxIcon /> : <MailIcon />}</ListItemIcon>
              <ListItemText primary={text} />
            </ListItem>
          ))}
        </List>
        <Divider />
        <List>
          {['All mail', 'Trash', 'Spam'].map((text, index) => (
            <ListItem button key={text}>
              <ListItemIcon>{index % 2 === 0 ? <InboxIcon /> : <MailIcon />}</ListItemIcon>
              <ListItemText primary={text} />
            </ListItem>
          ))}
        </List>
      </Drawer>
      <main className={classes.content}>
          { feed }
          {(props.loggedIn) &&
            <AddArticle newArticle={newArticle} setNewArticle={setNewArticle} />
          }
      </main>
    </div>
  
  if (error) {
    return (
      <div>
        Error: {error.message}
      </div>
    );
  } else if (!isLoaded) {
    return (
      <div>
        Loading...
      </div>
    );
  } else {
    return( 
      <div>
        {page} 
      </div>
    );
  }
}