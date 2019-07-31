import React, { useState, useEffect } from "react";
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import Avatar from '@material-ui/core/Avatar';
import mugshot from '../../images/mugshot.jpg';

const useStyles = makeStyles(theme => ({
    root: {
        display: 'flex',
        justifyContent: 'center',
    },
    paper: {
        justifyContent: 'center',
        width: '75%',
        padding: theme.spacing(3, 2)
    },
    typography: {
        marginTop: '5px',
        marginLeft: '5px',
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
                <Paper className={classes.paper} key={item.title}>
                    <Grid container spacing={1}>
                        <Grid item xs={12}>
                            <Typography variant="h3" gutterBottom>
                                {item.title}
                            </Typography>
                        </Grid>
                        <Grid container direction="row" justify="flex-start" alignItems="center">
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
                                <Typography variant="caption" gutterBottom>
                                    {item.category}
                                </Typography>
                            </Grid>
                            <Grid item>
                                <Typography variant="caption" gutterBottom>
                                    {item.topic}
                                </Typography>
                            </Grid>
                        </Grid>
                        <Grid item xs={12}>
                            <Typography component="p">
                                {item.content}
                            </Typography>
                        </Grid>
                    </Grid>
                </Paper>
            ));
        } else {
            setFeed(
                <div>Loading...</div>
            )
        }       
    }, [isLoaded]);

    return (
        (error) ? error.message : (!isLoaded) ? <div>Loading...</div> : <div className={classes.root}> {feed} </div>
    );
}