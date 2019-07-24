import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import Typography from '@material-ui/core/Typography';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import DirectionsRunIcon from '@material-ui/icons/DirectionsRun';
import TimerIcon from '@material-ui/icons/Timer';
import DevicesIcon from '@material-ui/icons/Devices';
import FitnessIcon from '@material-ui/icons/FitnessCenter';
import InboxIcon from '@material-ui/icons/MoveToInbox';
import Octicon, {Beaker} from '@primer/octicons-react';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardHeader from '@material-ui/core/CardHeader';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import AddArticle from '../AddArticle/AddArticle';
import FullArticle from '../FullArticle/FullArticle';

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
    height: 200,
    margin: theme.spacing(1),
    paddingBottom: 0
  },
  cardContent: {
    paddingTop: 0,
    // height: 0
  },
  cardActions: {
    margin: '3px'
  },
  typography: {
    marginBottom: '5px',
  },
  octicon: {
    height: '24px',
    width: '24px',
  },
  toolbar: theme.mixins.toolbar,
}));

export default function Articles(props) {
  const classes = useStyles();
  const [newArticle, setNewArticle] = useState(false);
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [articles, setArticles] = useState([]);
  const [feed, setFeed] = useState(null);
  const [filter, setFilter] = useState(null);

  const categories = ['Latest', 'Science', 'Technology', 'Sports', 'Health']
  const icons = [<TimerIcon />, <Octicon icon={Beaker} className={classes.octicon} />, <DevicesIcon />, <DirectionsRunIcon />, <FitnessIcon />]

  const handleFilter = text => () => {
    console.log(text);
    setFilter(text);
    fetch("http://172.23.90.20:8080/articles?cat=" + text.toLowerCase())
    .then(res => res.json())
    .then(result => {
      setArticles(result);
      setIsLoaded(true);
      setError(null);
    })
    .catch(error => {
      setIsLoaded(true);
      setError(error);
    })
  }

  const getResults = () => {
    //fetch("http://localhost:8080/home")
    fetch("http://172.23.90.20:8080/articles")
    .then(res => res.json())
    .then(result => {
        setArticles(result);
        setIsLoaded(true);
    })
    .catch(error => {
        setIsLoaded(true);
        setError(error);
    })
  };

  useEffect(() => {
    getResults();
  }, [newArticle]);

  const handleArticleClick = item => () => {
    setFeed(<FullArticle article={item} />)
  }

  useEffect(() => {
    if (articles.article) {
      console.log(articles);
      setFeed(articles.article.map(item =>
        <Card className={classes.card} key={item.title}>
          <CardActionArea onClick={handleArticleClick(item)}>
            <CardHeader
              marginBottom="5px"
              align="left" 
              title={item.title}
              subheader={item.author + ", " + item.date}
            />
            <CardContent className={classes.cardContent} align="left">
              <Typography className={classes.typography} variant="body2" color="textSecondary" component="p">
                {item.description}
              </Typography>
              <Divider />
              <Typography variant="caption" color="textSecondary" component="p">
                {item.category} / {item.topic}
              </Typography>
            </CardContent>
          </CardActionArea>
          <CardActions className={classes.cardActions}>
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
        <div></div>
      )
    }
  }, [articles, filter]);

  return (
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
          {categories.map((text, index) => (
            <ListItem button key={text} onClick={handleFilter(text)}>
              <ListItemIcon>
                {icons[index]}
              </ListItemIcon>
              <ListItemText primary={text} />
            </ListItem>
          ))}
          {(props.loggedIn) && 
            <ListItem button key="drafts" onClick={handleFilter("drafts")}>
              <ListItemIcon>
                <InboxIcon />
              </ListItemIcon>
              <ListItemText primary="Drafts" />
            </ListItem>
          }
        </List>
      </Drawer>
      <main className={classes.content}>
          { (error) ? error.message : (!isLoaded) ? <div>Loading...</div> : feed }
          {(props.loggedIn) &&
            <AddArticle newArticle={newArticle} setNewArticle={setNewArticle} />
          }
      </main>
    </div>
  );
}