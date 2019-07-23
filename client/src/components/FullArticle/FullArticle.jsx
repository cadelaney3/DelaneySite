import React, { useState, useEffect } from "react";
import { makeStyles } from '@material-ui/core/styles';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/List';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import Avatar from '@material-ui/core/Avatar';

const useStyles = makeStyles(theme => ({
    root: {
        display: 'flex',
        padding: theme.spacing(3, 2)
    },
    avatar: {
        margin: 10,
    },
}));

export default function FullArticle(props) {
    const classes = useStyles();

    return (
        <div>
            <Paper className={classes.root}>
                <Typography variant="h3" gutterBottom>
                    {props.article.title}
                </Typography>
                <Grid container alignItems="left">
                    <Avatar alt="Chris" src="../../images/mugshot.jpg" className={classes.avatar} />
                    <Typography variant="subtitle1" gutterBottom>
                        {props.article.author}
                    </Typography>
                    <Typography variant="subtitle1" gutterBottom>
                        {props.article.date}
                    </Typography>
                </Grid>
                <Grid container alignItems="left">
                    <Typography variant="caption" gutterBottom>
                        {props.article.category}
                    </Typography>
                    <Typography variant="caption" gutterBottom>
                        {props.article.topic}
                    </Typography>
                </Grid>
                <Typography variant="body1">
                    {props.article.content}
                </Typography>
            </Paper>
        </div>
    )
}