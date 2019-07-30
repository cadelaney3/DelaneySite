import React, { useState, useEffect } from "react";
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import Avatar from '@material-ui/core/Avatar';

const drawerWidth = 240;

const useStyles = makeStyles(theme => ({
    root: {
        display: 'flex',
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
    avatar: {
        margin: 10,
    },
    octicon: {
        height: '24px',
        width: '24px',
    },
}));

export default function FullArticle(props) {
    const classes = useStyles();
    const [error, setError] = useState(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [articles, setArticles] = useState([]);
    const [feed, setFeed] = useState(<div></div>);
    console.log(props);

    const getResult = () => {
        var query = "";
        if (props.location.article) {
            query = "?id=" + props.location.article.id;
        } else {
            query = "?title=" + window.location.pathname.split("/").pop();
        }
        fetch("http://localhost:8080/articles" + query)
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
        getResult();
    }, [])

    useEffect(() => {
        if (articles.article) {
            console.log("articles: ", articles);
            setFeed(articles.article.map(item =>
                <Paper className={classes.root}>
                    <Typography variant="h3" gutterBottom>
                        {item.title}
                    </Typography>
                    <Grid container alignItems="center">
                        <Avatar alt="Chris" src="../../images/mugshot.jpg" className={classes.avatar} />
                        <Typography variant="subtitle1" gutterBottom>
                            {item.author}
                        </Typography>
                        <Typography variant="subtitle1" gutterBottom>
                            {item.date}
                        </Typography>
                    </Grid>
                    <Grid container alignItems="center">
                        <Typography variant="caption" gutterBottom>
                            {item.category}
                        </Typography>
                        <Typography variant="caption" gutterBottom>
                            {item.topic}
                        </Typography>
                    </Grid>
                    <Typography component="p">
                        {item.content}
                    </Typography>
                </Paper>
            ));
        } else {
            setFeed(
                <div>Loading...</div>
            )
        }       
    }, [isLoaded]);

    return (
        (error) ? error.message : (!isLoaded) ? <div>Loading...</div> : feed
    );
}