import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { withRouter } from 'react-router';
import { Route, Link, NavLink, Redirect } from 'react-router-dom';
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

export default withRouter(function Articles(props) {
  const classes = useStyles();
  const [newArticle, setNewArticle] = useState(false);
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [articles, setArticles] = useState([]);
  const [feed, setFeed] = useState(null);
  const [filter, setFilter] = useState(null);
  const [article, setArticle] = useState(null);
  const [isArticleClicked, setIsArticleClicked] = useState(false);
  const [page, setPage] = useState(<div></div>);

  const categories = ['Latest', 'Science', 'Technology', 'Sports', 'Health']
  const icons = [<TimerIcon />, <Octicon icon={Beaker} className={classes.octicon} />, <DevicesIcon />, <DirectionsRunIcon />, <FitnessIcon />]

  const handleFilter = text => () => {
    setFilter(text);
  }

  const getResults = () => {
    console.log(filter);
    var query = "";
    if (filter !== null) {
      query = "?cat=" + filter.toLowerCase();
      console.log("filter is null");
    }
    fetch("http://172.25.59.60:8080/articles" + query)
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
  }, [newArticle, filter]);

  const handleArticleClick = item => () => {
    setArticle(<div><FullArticle article={item} /></div>);
    setIsArticleClicked(true);
  }

  useEffect(() => {
    if (articles.article) {
      console.log(articles);
      setFeed(articles.article.map(item =>
        <Link to={{pathname:`${props.match.url}/${item.title.replace(/\s+/g, '-').toLowerCase()}`, article: item}} key={item.id}>
        <Card className={classes.card} key={item.title}>
          <CardActionArea onClick={handleArticleClick(item)}> 
            <CardHeader
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
        </Link>
      ));
    } else {
      setFeed(
        <div></div>
      )
    }
  }, [articles, filter]);

  useEffect(() => {
    setPage(
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
          { (error) ? error.message : (!isLoaded) ? <div>Loading...</div> : feed } {/*(isArticleClicked) ? article : feed */}
          {(props.loggedIn) &&
            <AddArticle parentState={newArticle} setParentState={setNewArticle}/>
          }
      </main>
    </div>
    )
  }, [isArticleClicked, isLoaded, feed]);

  return (
    page
  );
})