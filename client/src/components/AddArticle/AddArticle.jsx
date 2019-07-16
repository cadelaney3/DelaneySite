import React, { useState, useEffect } from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import { makeStyles } from '@material-ui/core/styles';
import ListItemText from '@material-ui/core/ListItemText';
import ListItem from '@material-ui/core/ListItem';
import List from '@material-ui/core/List';
import Divider from '@material-ui/core/Divider';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import CloseIcon from '@material-ui/icons/Close';
import Slide from '@material-ui/core/Slide';


const headers = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Credentials": true,
  Accept: "text/plain"
};

const useStyles = makeStyles(theme => ({
    fab: {
        margin: theme.spacing(1),
        position: 'fixed',
        top: 'auto',
        bottom: 20,
        right: 20,
        left: 'auto',
    },
    form: {
        display: 'flex',
        flexDirection: 'column',
        margin: 'auto',
        width: 'fit-content',
    },
    formControl: {
    marginTop: theme.spacing(2),
    minWidth: 120,
    },
    formControlLabel: {
    marginTop: theme.spacing(1),
    },
    appBar: {
    position: 'relative',
    },
    title: {
    marginLeft: theme.spacing(2),
    flex: 1,
    },
    list: {
        display: 'flex',
        flexWrap: 'wrap',
    },
    textField: {
        marginLeft: theme.spacing(1),
        marginRight: theme.spacing(1),
    },
}));

const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
  });

export default function AddArticle(props) {
    const classes = useStyles();
    const [open, setOpen] = useState(false);
    const [article, setArticle] = useState({
        title: "",
        author: "",
        category: "",
        topic: "",
        description: "",
        content: "",
    });
  
    function handleClickOpen() {
      setOpen(true);
    }
  
    function handleClose() {
      setOpen(false);
      props.setNewArticle(false);
    }

    const handleChange = name => event => {
      //event.persist();
      //setArticle(event.target.value);
      setArticle({ ...article, [name]: event.target.value });
    }

    function handleSave() {
      fetch("http://172.26.34.14:8080/addArticle", {
        method: "POST",
        headers: headers,
        body: JSON.stringify({ title: article.title,
                               author: article.author,
                               category: article.category,
                               topic: article.topic,
                               description: article.description,
                               content: article.content
                            })
      })
        .then(results => results.json())
        .then(data => {
          if (data.status === 200) {
            props.setNewArticle(true);
            console.log(data);
          } else {
            console.log(data);
          }
        })
        .catch( err => {
          console.log(err);
          return Promise.reject();
        })
      handleClose();
    }
  
    return (
      <div>
        <Fab aria-label="Add" className={classes.fab}>
            <AddIcon onClick={handleClickOpen} />
        </Fab>
        <Dialog fullScreen open={open} onClose={handleClose} TransitionComponent={Transition}>
            <AppBar className={classes.appBar}>
                <Toolbar>
                    <IconButton edge="start" color="inherit" onClick={handleClose} aria-label="Close">
                        <CloseIcon />
                    </IconButton>
                    <Typography variant="h6" className={classes.title}>
                        Add New Article
                    </Typography>
                    <Button color="inherit" >
                        save as draft
                    </Button>
                    <Button color="inherit" onClick={handleSave}>
                        publish
                    </Button>
                </Toolbar>
            </AppBar>
            <List className={classes.list}>
                <ListItem>
                    <TextField
                        autoFocus
                        margin="dense"
                        id="title"
                        label="Title"
                        className={classes.textField}
                        type="text"
                        onChange={handleChange('title')}
                        multiline
                    />
                    <TextField
                        margin="dense"
                        id="author"
                        label="Author"
                        className={classes.textField}
                        type="text"
                        onChange={handleChange('author')}
                        multiline
                    />
                    <TextField
                        margin="dense"
                        id="category"
                        label="Category"
                        className={classes.textField}
                        type="text"
                        onChange={handleChange('category')}
                        multiline
                    />
                    <TextField
                        margin="dense"
                        id="topic"
                        label="Topic"
                        className={classes.textField}
                        type="text"
                        onChange={handleChange('topic')}
                        multiline
                    />
                </ListItem>
                <Divider />
                <ListItem>
                    <TextField
                        margin="dense"
                        id="description"
                        label="Description"
                        className={classes.textField}
                        type="text"
                        onChange={handleChange('description')}
                        variant="outlined"
                        rows="3"
                        multiline
                        fullWidth
                    />
                </ListItem>
                <Divider />
                <ListItem>
                    <TextField
                        // autoFocus
                        margin="normal"
                        id="content"
                        label="Content"
                        className={classes.textField}
                        type="text"
                        onChange={handleChange('content')}
                        variant="outlined"
                        rows="30"
                        multiline
                        fullWidth
                    />
                </ListItem>
                <Divider />
                <ListItem button>
                    <ListItemText primary="Default notification ringtone" secondary="Tethys" />
                </ListItem>
            </List>
        </Dialog>
      </div>
    );
}