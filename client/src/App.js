import React, { Component, useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import './App.css';
import Header from './components/Header/Header';
import ChatHistory from './components/ChatHistory/ChatHistory';
import ChatInput from './components/ChatInput/ChatInput';
import Home from './components/Home/Home';
import SignIn from './components/SignIn/SignIn';
import Articles from './components/Articles/Articles';
import { connect, sendMsg } from "./api";
import { createMuiTheme } from '@material-ui/core/styles';
import { MuiThemeProvider } from '@material-ui/core/styles';

const theme = createMuiTheme({
  palette: {
    type: 'dark',
  }
})

export default function App() {
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    sessionStorage.setItem('loggedIn', loggedIn);
  });

  const handleSignInChange = () => {
    setLoggedIn(!loggedIn);
  }

  return (
    <Router>
      <MuiThemeProvider theme={theme} >
        <Header theme={theme} loggedIn={loggedIn} handleSignInChange={handleSignInChange} />
        <Switch>
          <Route exact path="/ws" render={(props) => <WS {...props} loggedIn={loggedIn} /> } />
          <Route 
            exact path="/articles/" 
            render={(props) => <Articles
              {...props}
              loggedIn={loggedIn}
              />
            }
          />
          <Route 
            exact path="/signin" 
            render={(props) => <SignIn
              {...props}
              loggedIn={loggedIn} 
              handleSignInChange={handleSignInChange}
              />            
            } 
          />
          <Route 
            expact path="/" 
            render={(props) => <Home 
              {...props} 
              loggedIn={loggedIn} 
              />
            }
          />
        </Switch>
      </MuiThemeProvider>
    </Router>
  );
}

class WS extends Component {
  constructor(props) {
    super(props);

    this.state = {
      prevState: null,
      chatHistory: []
    }
  }

  componentDidMount() {
    connect((msg) => {
      console.log("New Message")
      this.setState(prevState => ({
        chatHistory: [...this.state.chatHistory, msg]
      }))
      console.log(this.state);
    });
  }

  send(event) {
    if(event.keyCode === 13) {
      sendMsg(event.target.value);
      event.target.value = "";
    }
  }

  render() {
    return (
      <div className="WS">
        <ChatHistory chatHistory={this.state.chatHistory} />
        <ChatInput send={this.send} />
      </div>
    );
  }
}
