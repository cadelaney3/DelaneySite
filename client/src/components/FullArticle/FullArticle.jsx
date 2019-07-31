import React, { useState, useEffect } from "react";
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import mugshot from '../../images/mugshot.jpg';
import AddArticle from '../AddArticle/AddArticle';

const useStyles = makeStyles(theme => ({
    root: {
        display: 'flex',
        justifyContent: 'center',
        padding: theme.spacing(4)
    },
    paper: {
        justifyContent: 'center',
        width: '75%',
        minHeight: '100vw',
        padding: theme.spacing(3, 8)
    },
    typography: {
        marginTop: '5px',
        marginLeft: '5px',
    },
    contentGrid: {
        margin: theme.spacing(8,0)
    },
    button: {
        margin: theme.spacing(1),
    }
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
        getResult();
    }, [])

    // const handleEdit = () => {
    //     if (articles.article) {
    //         <AddArticle content={articles.article} />
    //     }
    // }

    useEffect(() => {
        if (articles.article) {
            // console.log("articles: ", articles);
            setFeed(articles.article.map(item =>
                <div key={item.title} className={classes.root}>
                {(props.loggedIn) &&
                    <AddArticle content={item} />
                }               
                <Paper className={classes.paper}>
                    <Grid container spacing={1}>
                        <Grid item xs={12}>
                            <Typography variant="h2" gutterBottom>
                                {item.title}
                            </Typography>
                        </Grid>
                        <Grid container direction="row" justify="flex-start" alignItems="center" spacing={1}>
                            <Grid item>
                                <Avatar alt="Chris" src={mugshot} className={classes.avatar} />
                            </Grid>
                            <Grid item>
                                <Typography className={classes.typography} variant="subtitle1" gutterBottom>
                                    {item.author}
                                </Typography>
                            </Grid>
                            <Grid item>
                                <Typography className={classes.typography} variant="subtitle1" gutterBottom>
                                    {item.date}
                                </Typography>
                            </Grid>
                        </Grid>
                        <Grid container direction="row" justify="flex-start" alignItems="center">
                            <Grid item>
                                <Typography className={classes.typography} variant="caption" gutterBottom>
                                    {item.category}
                                </Typography>
                            </Grid>
                            <Grid item>
                                <Typography className={classes.typography} variant="caption" gutterBottom>
                                    {item.topic}
                                </Typography>
                            </Grid>
                        </Grid>
                        <Grid className={classes.contentGrid} item xs={12}>
                            <Typography component="p">
                                {item.content}
                            </Typography>
                        </Grid>
                    </Grid>
                </Paper>
                </div>
            ));
        } else {
            setFeed(
                <div>Loading...</div>
            )
        }       
    }, [isLoaded]);

    //return (
        // (error) ? error.message : (!isLoaded) ? <div>Loading...</div> : <div className={classes.root}> {feed} </div>
    if (error) {
        return (<div>Error: {error.message}</div>);
    } else if (!isLoaded) {
        return (<div>Loading...</div>);
    } else {
        return (
            feed
        );
    }
    //);
}