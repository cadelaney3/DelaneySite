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
    }
}));

export default function AddFact(props) {
    const classes = useStyles();
    const [open, setOpen] = useState(false);
    const [fact, setFact] = useState("")
  
    function handleClickOpen() {
      setOpen(true);
    }
  
    function handleClose() {
      setOpen(false);
      props.setNewFact(false);
    }

    function handleFactChange(event) {
      event.persist();
      setFact(event.target.value);
    }

    function handleSave() {
      fetch("http://172.17.21.104:8080/addFact", {
        method: "POST",
        headers: headers,
        body: JSON.stringify({ fact: fact })
      })
        .then(results => results.json())
        .then(data => {
          if (data.status === 200) {
            props.setNewFact(true);
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
        <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
          <DialogTitle id="form-dialog-title">Add New Fact</DialogTitle>
          <DialogContent>
            <DialogContentText>
              Add a new fact by typing in the text field below.
            </DialogContentText>
            <TextField
              autoFocus
              margin="dense"
              id="fact"
              label="Fact"
              type="text"
              onChange={handleFactChange}
              multiline
              fullWidth
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">
              Cancel
            </Button>
            <Button onClick={handleSave} color="primary">
              Save
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }

