import React, { useState, useEffect } from 'react';
import clsx from 'clsx';
import { makeStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import CssBaseline from '@material-ui/core/CssBaseline';
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
import CardMedia from '@material-ui/core/CardMedia';
import Button from '@material-ui/core/Button';
import lakecomo from '../../images/lakecomo.jpg';

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
    justifyContent: 'space-evenly',
    flexGrow: 1,
    padding: theme.spacing(3),
  },
  card: {
    width: 450,
    //margin: theme.spacing(1)
  },
  toolbar: theme.mixins.toolbar,
}));

function ImgMediaCard() {
  const classes = useStyles();
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [articles, setArticles] = useState([]);
  const [newArticle, setNewArticle] = useState(false);
  const [feed, setFeed] = useState([]);

  const getResults = () => {
    //fetch("http://localhost:8080/home")
    fetch("http://172.17.251.115:8080/articles")
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
  }, [newArticle]);

  useEffect(() => {
    if (isLoaded) {
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
    } else {
      setFeed(
        <div>
          Loading...
        </div>
      )
    }
  }, [isLoaded]);

  return (
    feed
    // <Card className={classes.card}>
    //   <CardActionArea>
    //     <CardMedia
    //       component="img"
    //       alt="Contemplative Reptile"
    //       height="140"
    //       image={lakecomo}
    //       title="Contemplative Reptile"
    //     />
    //     <CardContent align="left">
    //       <Typography gutterBottom variant="h5" component="h2">
    //         Lizard
    //       </Typography>
    //       <Typography variant="body2" color="textSecondary" component="p">
    //         Lizards are a widespread group of squamate reptiles, with over 6,000 species, ranging
    //         across all continents except Antarctica
    //       </Typography>
    //     </CardContent>
    //   </CardActionArea>
    //   <CardActions>
    //     <Button size="small" color="primary">
    //       Share
    //     </Button>
    //     <Button size="small" color="primary">
    //       Learn More
    //     </Button>
    //   </CardActions>
    // </Card>
  );
}

export default function ClippedDrawer() {
  const classes = useStyles();

  return (
    <div className={classes.root} align='center'>
      {/* <CssBaseline />
      <AppBar position="fixed" className={classes.appBar}>
        <Toolbar>
          <Typography variant="h6" noWrap>
            Clipped drawer
          </Typography>
        </Toolbar>
      </AppBar> */}
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
        {/* <div className={classes.toolbar}/> */}
          {ImgMediaCard()}
      </main>
    </div>
  );
}